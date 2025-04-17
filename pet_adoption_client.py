import grpc
import pet_adoption_pb2
import pet_adoption_pb2_grpc

def register_pet(stub):
    name = input("Enter the pet's name: ")
    gender = input("Enter the pet's gender (Male/Female): ")
    age_input = input("Enter the pet's age: ")
    breed = input("Enter the pet's breed: ")
    picture = input("Enter the path to the pet's picture: ")

    try:
        age = int(age_input)
    except ValueError:
        print("Invalid age input. Please enter a valid number.")
        return  # Return to menu

    # Register a new pet with the given details
    response = stub.RegisterPet(
        pet_adoption_pb2.RegisterPetRequest(
            name=name, gender=gender, age=age, breed=breed, picture=picture
        )
    )

    if "Error:" in response.message:
        print(response.message)
        return  # Return to menu if there was an error

    print(f"RegisterPet Response: {response.message}")

def search_pet(stub):
    name = input("Enter the pet's name to search (or leave blank): ")
    gender = input("Enter the pet's gender to search (or leave blank): ")
    age_input = input("Enter the pet's age to search (or leave blank): ")
    age = int(age_input) if age_input else 0
    breed = input("Enter the pet's breed to search (or leave blank): ")

    search_response = stub.SearchPet(
        pet_adoption_pb2.SearchPetRequest(
            name=name, gender=gender, age=age, breed=breed
        )
    )

    print(f"Found {len(search_response.pets)} matching pet(s):")
    for pet in search_response.pets:
        print(f"Name: {pet.name}, Breed: {pet.breed}, Age: {pet.age}, Gender: {pet.gender}, Picture: {pet.picture}")

def run():
    with grpc.insecure_channel('host.docker.internal:50051') as channel:
        stub = pet_adoption_pb2_grpc.PetAdoptionServiceStub(channel)

        while True:
            print("\nPet Adoption System")
            print("1. Register a new pet")
            print("2. Search for a pet")
            print("3. Exit")
            choice = input("Choose an option: ")

            if choice == '1':
                register_pet(stub)
            elif choice == '2':
                search_pet(stub)
            elif choice == '3':
                print("Exiting...")
                break
            else:
                print("Invalid choice, please try again.")

if __name__ == '__main__':
    run()
