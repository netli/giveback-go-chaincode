package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var MarketplaceStr = "_marketplace"       // name for the key/value that will store all open tasks
var CompletedTasksStr = "_completedTasks" // name for the key/value that will store all completed tasks (all tasks = marketplace + completedtasks)

type Task struct {
	Uid         string   `json:"id"`
	User        string   `json:"user"` // users are defined by their emails
	Amount      int      `json:"amount"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Skills      string   `json:"skills"`
	Location    string   `json:"location"` // either "remote" or "onsite"
	Address     string   `json:"address"`
	Submissions []string `json:"submissions"`
	CompletedBy string   `json:"completedBy"` // either null or a user's email
}

type Marketplace struct {
	Tasks []Task `json:"openTasks"`
}

type CompletedTasks struct { // all tasks here shoud have Completed by not null
	Tasks []Task `json:"closedTasks"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval))) //making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	var cTasks CompletedTasks
	jsonAsBytes, _ := json.Marshal(cTasks) //clear the CompletedTasks struct
	err = stub.PutState(CompletedTasksStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	var mplace Marketplace
	jsonAsBytes, _ = json.Marshal(mplace) //clear the Marketplace struct
	err = stub.PutState(MarketplaceStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Invoke - Our entry point to invoke a chaincode function (eg. write, createAccount, etc)
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "add_task" {
		return t.add_task(stub, args)
	} else if function == "add_submission" {
		return t.add_submission(stub, args)
	} else if function == "delete_submission" {
		return t.delete_submission(stub, args)
	} else if function == "end_task" {
		return t.end_task(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// ============================================================================================================================
// Read - read a variable from chaincode state (used by Query)
// ============================================================================================================================

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// ============================================================================================================================
// Write - Invoke function to write
// ============================================================================================================================
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// add_task - creates a task and adds it to Marketplace struct
// ============================================================================================================================
func (t *SimpleChaincode) add_task(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2         3           4              5           6          7          8          9
	// "uid", "user", "amount", "title", "description", "start date", "end date", "skills", "location", "address" (address is optional)
	if len(args) < 9 {
		return nil, errors.New("Incorrect number of arguments. Expecting 9 or 10")
	}

	fmt.Println("- create and add task")

	fmt.Println("Number of args: ")
	fmt.Println(len(args))
	fmt.Println(args[0])
	fmt.Println(args[1])
	fmt.Println(args[2])
	fmt.Println(args[3])
	fmt.Println(args[4])
	fmt.Println(args[5])
	fmt.Println(args[6])
	fmt.Println(args[7])
	fmt.Println(args[8])

	amount, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument (amount) must be a numeric string")
	}

	var task = Task{}
	task.Uid = args[0]
	task.User = args[1]
	task.Amount = amount
	task.Title = args[3]
	task.Description = args[4]
	task.StartDate = args[5]
	task.EndDate = args[6]
	task.Skills = args[7]
	task.Location = args[8]
	if len(args) > 9 {
		fmt.Println(args[9])
		fmt.Println("address not empty, add to task")
		task.Address = args[9]
	}

	fmt.Println("below is task: ")
	fmt.Println(task)

	////////////////////// 1) store task with Uid as key for easy search /////
	taskAsBytes, _ := json.Marshal(task)      /////
	err = stub.PutState(args[0], taskAsBytes) /////
	if err != nil {                           /////
		return nil, err /////
	} /////
	//////////////////////////////////////////////////////////////////////////

	//get the marketplace struct
	MarketplaceAsBytes, err := stub.GetState(MarketplaceStr)
	if err != nil {
		return nil, errors.New("Failed to get marketplace")
	}
	var mplace Marketplace
	json.Unmarshal(MarketplaceAsBytes, &mplace) //un stringify it aka JSON.parse()

	/////////////////////// 2) append task into marketplace ///////////////////////////
	mplace.Tasks = append(mplace.Tasks, task)        /////
	fmt.Println("! appended task to marketplace")    /////
	jsonAsBytes, _ := json.Marshal(mplace)           /////
	err = stub.PutState(MarketplaceStr, jsonAsBytes) //rewrite marketplace		  /////
	if err != nil {                                  /////
		return nil, err /////
	} /////
	///////////////////////////////////////////////////////////////////////////////////

	fmt.Println("- end add task")
	return nil, nil
}

// ============================================================================================================================
// modify_task - update/delete submission on task, or add completedBy. Only for single task key, not marketplace/CompletedTasks
// ============================================================================================================================
func (t *SimpleChaincode) modify_task(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0          1           2
	// "command", "uid",  "bob@email.com" (email is optional)
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2 or 3")
	}

	fmt.Println("- start helper modify_task")
	// fmt.Println(args[0] + " - " + args[1] + " - " + args[2])

	// get task from blockchain
	tasksAsBytes, err := stub.GetState(args[1])
	if err != nil {
		return nil, errors.New("Failed to get task")
	}

	res := Task{}
	json.Unmarshal(tasksAsBytes, &res) //un stringify it aka JSON.parse()

	if args[0] == "add_submission" {
		res.Submissions = append(res.Submissions, args[2]) // append submission
		fmt.Println("! appended submission to task")
		fmt.Println(res.Submissions)
		fmt.Println(res)
	} else if args[0] == "delete_submission" {
		for i, v := range res.Submissions { // remove submission from task
			if v == args[2] {
				fmt.Println("found v")
				res.Submissions = append(res.Submissions[:i], res.Submissions[i+1:]...)
				break
			}
		}
		fmt.Println("! deleted submission from task")
		fmt.Println(res.Submissions)
		fmt.Println(res)
	} else if args[0] == "add_completedBy" {
		if len(args) == 3 {
			res.CompletedBy = args[2]
			fmt.Println("! marked task as complete by user")
		} else {
			res.CompletedBy = "CLOSED"
			fmt.Println("! marked task as CLOSED")
		}
		fmt.Println(res.CompletedBy)
		fmt.Println(res)
	} else {
		fmt.Println("======UNKNOWN OPERATION======")
	}

	// update task and push back into blockchain
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[1], jsonAsBytes) //rewrite the task with id as key
	if err != nil {
		return nil, err
	}

	fmt.Println("- end helper modify_task")
	return nil, nil
}

// ============================================================================================================================
// add_submission - upload submission into task
// ============================================================================================================================
func (t *SimpleChaincode) add_submission(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Number of args: ")
	fmt.Println(len(args))
	fmt.Println(args[0])
	fmt.Println(args[1])

	//   0            1
	// "uid", "user_who_updated"
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2: task uid and user who updated")
	}

	fmt.Println("- add submission")

	// get all active tasks in marketplace
	MarketplaceAsBytes, err := stub.GetState(MarketplaceStr)
	if err != nil {
		return nil, errors.New("Failed to get marketplace array")
	}
	var mplace Marketplace
	json.Unmarshal(MarketplaceAsBytes, &mplace) //un stringify it aka JSON.parse()

	fmt.Print("Marketplace array: ")
	fmt.Println(mplace)

	//////// update submission in task in marketplace //////
	for i := range mplace.Tasks { //iter through all the tasks
		fmt.Print("looking @ task name: ")
		fmt.Println(mplace.Tasks[i])

		if mplace.Tasks[i].Uid == args[0] { // found the trade to update
			fmt.Println("Found trade to add submission")

			t.modify_task(stub, []string{"add_submission", args[0], args[1]}) // add submission to single uid query

			mplace.Tasks[i].Submissions = append(mplace.Tasks[i].Submissions, args[1]) // add submission to marketplace array
			fmt.Println("! appended submission to task in marketplace")
			fmt.Println(mplace.Tasks[i].Submissions)
			fmt.Println(mplace.Tasks[i])

			jsonAsBytes, _ := json.Marshal(mplace)
			err = stub.PutState(MarketplaceStr, jsonAsBytes) //rewrite the marketplace with new submission
			if err != nil {
				return nil, err
			}
			break
		} else if i == (len(mplace.Tasks) - 1) {
			return nil, errors.New("!Task not found in add_submission")
		}
	}
	fmt.Println("- end add submission")
	return nil, nil
}

// ============================================================================================================================
// delete_submission - delete submission from task
// ============================================================================================================================
func (t *SimpleChaincode) delete_submission(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Number of args: ")
	fmt.Println(len(args))
	fmt.Println(args[0])
	fmt.Println(args[1])

	//   0            1
	// "uid", "user_who_updated"
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2: task uid and user who updated")
	}

	fmt.Println("- delete submission")

	// get all active tasks in marketplace
	MarketplaceAsBytes, err := stub.GetState(MarketplaceStr)
	if err != nil {
		return nil, errors.New("Failed to get marketplace array")
	}
	var mplace Marketplace
	json.Unmarshal(MarketplaceAsBytes, &mplace) //un stringify it aka JSON.parse()

	fmt.Print("Marketplace array: ")
	fmt.Println(mplace)

	//////// update submission in task in marketplace //////
	for i := range mplace.Tasks { //iter through all the tasks
		fmt.Print("looking @ task name: ")
		fmt.Println(mplace.Tasks[i])

		if mplace.Tasks[i].Uid == args[0] { // found the trade to update
			fmt.Println("Found trade to delete submission")

			if len(mplace.Tasks[i].Submissions) == 0 { // if task has no submissions
				return nil, errors.New("!Did not find corresponding submission in task to delete")
			}

			for j, v := range mplace.Tasks[i].Submissions { // delete submission from marketplace array
				fmt.Println(mplace.Tasks[i].Submissions)
				if v == args[1] {
					fmt.Println("found v")
					t.modify_task(stub, []string{"delete_submission", args[0], args[1]}) // delete submission from single uid query
					mplace.Tasks[i].Submissions = append(mplace.Tasks[i].Submissions[:j], mplace.Tasks[i].Submissions[j+1:]...)
					break
				} else if j == (len(mplace.Tasks[i].Submissions) - 1) { // did not find submission to delete
					return nil, errors.New("!Did not find corresponding submission in task to delete")
				}
			}

			fmt.Println("! deleted submission from task in marketplace")
			fmt.Println(mplace.Tasks[i].Submissions)
			fmt.Println(mplace.Tasks[i])

			jsonAsBytes, _ := json.Marshal(mplace)
			err = stub.PutState(MarketplaceStr, jsonAsBytes) //rewrite the marketplace with new submission
			if err != nil {
				return nil, err
			}
			break
		} else if i == (len(mplace.Tasks) - 1) {
			return nil, errors.New("!Task not found in delete_submission")
		}
	}
	fmt.Println("- end delete submission")
	return nil, nil
}

// ============================================================================================================================
// end_task - (if user_who_finished is given, set task as finished by user and) move task to CompletedTrade
// ============================================================================================================================
func (t *SimpleChaincode) end_task(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0            1
	// "uid", "user_who_finished" (user_who_finished is optional)
	if len(args) > 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 or 2")
	}

	fmt.Println("- start end task")

	// get all active tasks in marketplace
	MarketplaceAsBytes, err := stub.GetState(MarketplaceStr)
	if err != nil {
		return nil, errors.New("Failed to get marketplace array")
	}
	var mplace Marketplace
	json.Unmarshal(MarketplaceAsBytes, &mplace) //un stringify it aka JSON.parse()

	fmt.Print("Marketplace array: ")
	fmt.Println(mplace)

	var completedTask = Task{}

	//////// update completedBy in task in marketplace //////
	for i := range mplace.Tasks { //iter through all the tasks
		fmt.Print("looking @ task name: ")
		fmt.Println(mplace.Tasks[i])

		if mplace.Tasks[i].Uid == args[0] { // found the trade to update
			fmt.Println("Found trade to delete in marketplace array")

			if len(args) == 2 { // if task is finished by a user
				t.modify_task(stub, []string{"add_completedBy", args[0], args[1]})
				mplace.Tasks[i].CompletedBy = args[1] // add user to completedBy
			} else {
				t.modify_task(stub, []string{"add_completedBy", args[0]})
				mplace.Tasks[i].CompletedBy = "CLOSED"
			}

			completedTask = mplace.Tasks[i]
			fmt.Println(completedTask)
			mplace.Tasks = append(mplace.Tasks[:i], mplace.Tasks[i+1:]...) // remove task from marketplace
			fmt.Println(mplace)

			jsonAsBytes, _ := json.Marshal(mplace)
			err = stub.PutState(MarketplaceStr, jsonAsBytes) //rewrite the marketplace with new submission
			if err != nil {
				return nil, err
			}
			break
		} else if i == (len(mplace.Tasks) - 1) {
			return nil, errors.New("!Task not found in end_task")
		}
	}

	//get the Completed Tasks struct
	CompletedTasksAsBytes, err := stub.GetState(CompletedTasksStr)
	if err != nil {
		return nil, errors.New("Failed to get marketplace")
	}
	var cTasks CompletedTasks
	json.Unmarshal(CompletedTasksAsBytes, &cTasks) //un stringify it aka JSON.parse()

	//// append task into marketplace ///////////////////////////
	cTasks.Tasks = append(cTasks.Tasks, completedTask)
	fmt.Println("! appended task to CompletedTasks")
	jsonAsBytes, _ := json.Marshal(cTasks)
	err = stub.PutState(CompletedTasksStr, jsonAsBytes) //rewrite marketplace
	if err != nil {
		return nil, err
	}

	fmt.Print("new marketplace: ")
	fmt.Println(mplace)
	fmt.Print("new completedTasks: ")
	fmt.Println(cTasks)

	fmt.Println("- end end task")

	return nil, nil
}
