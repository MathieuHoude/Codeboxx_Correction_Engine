require 'residential_controller'

def deep_copy(o)
    Marshal.load(Marshal.dump(o))
  end

def scenario(column,requestedFloor, direction, destination)
    tempColumn = deep_copy(column)
    selectedElevator = tempColumn.requestElevator(requestedFloor, direction).clone
    pickedUpUser = true if selectedElevator.floor == requestedFloor
    selectedElevator.requestFloor(destination, tempColumn)
    tempColumn.elevatorsList.select {|elevator| elevator.id == selectedElevator.id}[0].floor = selectedElevator.floor
    

    results = {
        :tempColumn => tempColumn,
        :selectedElevator => selectedElevator,
        :pickedUpUser => pickedUpUser
    }

end

describe ResidentialController do
    describe "Column's attributes and methods" do
        column = ResidentialController::Column.new(1,'active',10,2)

        it 'instantiates a Column with valid attributes' do 
            expect(column).to be_a(ResidentialController::Column)
            .and have_attributes(
                :id => 1,
                :status => 'active',
                :numberOfFloors => 10,
                :numberOfElevators => 2,
                :elevatorsList => array_including(kind_of(ResidentialController::Elevator)),
                :buttonsUpList => array_including(kind_of(ResidentialController::Button)),
                :buttonsDownList => array_including(kind_of(ResidentialController::Button))
            )
            expect(column.elevatorsList.count).to eq(2)
            expect(column.buttonsUpList.count).to eq(9)
            expect(column.buttonsDownList.count).to eq(9)
        end

        it 'has a requestElevator method' do
            expect(column).to respond_to(:requestElevator)
        end

        it 'can find and return an elevator' do
            expect(column.findElevator(1, 'up')).to be_an(ResidentialController::Elevator)
        end
    end

    describe "Elevator's attributes and methods" do
        elevator = ResidentialController::Elevator.new(1, 10, 1, 'idle', 'off', 'off' )

        it 'instantiates an Elevator with valid attributes' do 
            expect(elevator).to be_a(ResidentialController::Elevator)
            .and have_attributes(
                :id => 1,
                :numberOfFloors => 10,
                :floor => 1,
                :status => 'idle',
                :weightSensor => 'off',
                :obstructionSensor => 'off',
                :elevatorDoor => be_a(ResidentialController::Door),
                :floorButtonsList => array_including(kind_of(ResidentialController::Button)),
                :floorRequestList => []
            )
            expect(elevator.floorButtonsList.count).to eq(10)
        end

        it 'has a requestFloor method' do
            expect(elevator).to respond_to(:requestFloor)
        end

    end

    describe "Functionnal scenario 1" do 
        column = ResidentialController::Column.new(1,'active',10,2)
        column.elevatorsList[0].floor = 2
        column.elevatorsList[1].floor = 6

        results = scenario(column, 3, 'up', 7)

        it "chooses the best elevator" do
            expect(results[:selectedElevator].id).to eq(1)
        end

        it "picks up the user at it's floor" do
            expect(results[:pickedUpUser]).to be(true)
        end
        
        it "brings the user to it's destination" do
            expect(results[:selectedElevator].floor).to eq(7)
        end

        it "ends with all the elevators at the right position" do
            expect(results[:tempColumn].elevatorsList[0].floor).to eq(7)
            expect(results[:tempColumn].elevatorsList[1].floor).to eq(6)
        end
    end

    describe "Functionnal scenario 2" do 
        column = ResidentialController::Column.new(1,'active',10,2)
        column.elevatorsList[0].floor = 10
        column.elevatorsList[1].floor = 3

        results1 = scenario(column, 1, 'up', 6)

        column = results1[:tempColumn] # update the column state with last scenario's result
        
        results2 = scenario(column, 3, 'up', 5)
        column = results2[:tempColumn] # update the column state with last scenario's result

        results3 = scenario(column, 9, 'down', 2)
        column = results3[:tempColumn] # update the column state with last scenario's result

        context "Part 1 of scenario 2" do
            it "chooses the best elevator" do
                expect(results1[:selectedElevator].id).to eq(2)
            end

            it "picks up the user at it's floor" do
                expect(results1[:pickedUpUser]).to be(true)
            end
            
            it "brings the user to it's destination" do
                expect(results1[:selectedElevator].floor).to eq(6)
            end

            it "ends with all the elevators at the right position" do
                expect(results1[:tempColumn].elevatorsList[0].floor).to eq(10)
                expect(results1[:tempColumn].elevatorsList[1].floor).to eq(6)
            end
        end

        context "Part 2 of scenario 2" do
            it "chooses the best elevator" do
                expect(results2[:selectedElevator].id).to eq(2)
            end

            it "picks up the user at it's floor" do
                expect(results2[:pickedUpUser]).to be(true)
            end
            
            it "brings the user to it's destination" do
                expect(results2[:selectedElevator].floor).to eq(5)
            end

            it "ends with all the elevators at the right position" do
                expect(results2[:tempColumn].elevatorsList[0].floor).to eq(10)
                expect(results2[:tempColumn].elevatorsList[1].floor).to eq(5)
            end
        end

        context "Part 3 of scenario 2" do
            it "chooses the best elevator" do
                expect(results3[:selectedElevator].id).to eq(1)
            end

            it "picks up the user at it's floor" do
                expect(results3[:pickedUpUser]).to be(true)
            end
            
            it "brings the user to it's destination" do
                expect(results3[:selectedElevator].floor).to eq(2)
            end

            it "ends with all the elevators at the right position" do
                expect(results3[:tempColumn].elevatorsList[0].floor).to eq(2)
                expect(results3[:tempColumn].elevatorsList[1].floor).to eq(5)
            end
        end
    end

    describe "Functionnal scenario 3" do 
        column = ResidentialController::Column.new(1,'active',10,2)
        column.elevatorsList[0].floor = 10
        column.elevatorsList[1].floor = 3
        column.elevatorsList[1].status = 'movingUp'

        results1 = scenario(column, 3, 'down', 2)

        column = results1[:tempColumn] # update the column state with last scenario's result

        column.elevatorsList[1].floor = 6
        column.elevatorsList[1].status = 'idle'

        results2 = scenario(column, 10, 'down', 3)
        column = results2[:tempColumn] # update the column state with last scenario's result

        context "Part 1 of scenario 3" do
            it "chooses the best elevator" do
                expect(results1[:selectedElevator].id).to eq(1)
            end

            it "picks up the user at it's floor" do
                expect(results1[:pickedUpUser]).to be(true)
            end
            
            it "brings the user to it's destination" do
                expect(results1[:selectedElevator].floor).to eq(2)
            end

            it "ends with all the elevators at the right position" do
                expect(results1[:tempColumn].elevatorsList[0].floor).to eq(2)
                expect(results1[:tempColumn].elevatorsList[1].floor).to eq(6)
            end
        end

        context "Part 2 of scenario 3" do
            it "chooses the best elevator" do
                expect(results2[:selectedElevator].id).to eq(2)
            end

            it "picks up the user at it's floor" do
                expect(results2[:pickedUpUser]).to be(true)
            end
            
            it "brings the user to it's destination" do
                expect(results2[:selectedElevator].floor).to eq(3)
            end

            it "ends with all the elevators at the right position" do
                expect(results2[:tempColumn].elevatorsList[0].floor).to eq(2)
                expect(results2[:tempColumn].elevatorsList[1].floor).to eq(3)
            end
        end
    end
end