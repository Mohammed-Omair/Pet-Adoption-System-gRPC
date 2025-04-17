package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
	pb "pet_adoption/pb"
)

type server struct {
	pb.UnimplementedPetAdoptionServiceServer
	mu     sync.Mutex
	petDB  string
}

type Pet struct {
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Age     int32  `json:"age"`
	Breed   string `json:"breed"`
	Picture string `json:"picture"`
}

// Initialize the server and load the pet database file
func (s *server) loadPets() ([]Pet, error) {
	var pets []Pet
	file, err := os.Open(s.petDB)
	if err != nil {
		if os.IsNotExist(err) {
			return pets, nil // No file yet, return empty list
		}
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&pets)
	return pets, err
}

// Save the pet data to JSON file
func (s *server) savePets(pets []Pet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create(s.petDB)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(&pets)
	return err
}

// Register a new pet
func (s *server) RegisterPet(ctx context.Context, req *pb.RegisterPetRequest) (*pb.RegisterPetResponse, error) {
	// Validation: Check for missing details and negative age
	if req.GetName() == "" || req.GetGender() == "" || req.GetAge() <= 0 || req.GetBreed() == "" || req.GetPicture() == "" {
		return &pb.RegisterPetResponse{Message: "Error: Required details not entered or invalid age."}, nil
	}

	pets, err := s.loadPets()
	if err != nil {
		return nil, err
	}

	// Add new pet to the list
	newPet := Pet{
		Name:    req.GetName(),
		Gender:  req.GetGender(),
		Age:     req.GetAge(),
		Breed:   req.GetBreed(),
		Picture: req.GetPicture(),
	}
	pets = append(pets, newPet)

	// Save to file
	if err := s.savePets(pets); err != nil {
		return nil, err
	}

	return &pb.RegisterPetResponse{Message: "Pet registered successfully!"}, nil
}

// Search for a pet
func (s *server) SearchPet(ctx context.Context, req *pb.SearchPetRequest) (*pb.SearchPetResponse, error) {
	pets, err := s.loadPets()
	if err != nil {
		return nil, err
	}

	var matchedPets []*pb.Pet
	for _, pet := range pets {
		if (req.GetName() == "" || req.GetName() == pet.Name) &&
			(req.GetGender() == "" || req.GetGender() == pet.Gender) &&
			(req.GetAge() == 0 || req.GetAge() == pet.Age) &&
			(req.GetBreed() == "" || req.GetBreed() == pet.Breed) {
			matchedPets = append(matchedPets, &pb.Pet{
				Name:    pet.Name,
				Gender:  pet.Gender,
				Age:     pet.Age,
				Breed:   pet.Breed,
				Picture: pet.Picture,
			})
		}
	}

	return &pb.SearchPetResponse{Pets: matchedPets}, nil
}

func main() {
	// Initialize server with the JSON file to store pet data
	s := &server{petDB: "/app/data/pets.json"}

	// Create a TCP listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the PetAdoptionService on the server
	pb.RegisterPetAdoptionServiceServer(grpcServer, s)

	fmt.Println("Go gRPC server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
