package helpers

import(
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)


//TODO: currently not in used, but this is a ready function to be tested

var encoderScript = `import numpy as np
import pandas as pd
import torch
import torch.nn as nn
import torch.nn.parallel
import torch.optim as optim
import torch.utils.data
from torch.autograd import Variable


# Preparing the training set and the test set
training_set = %v
training_set = np.array(training_set, dtype = 'float')
test_set = %v
test_set = np.array(test_set, dtype = 'float')
predict_set = %v
predict_set = np.array(predict_set, dtype = 'float')

# Getting the number of attributes and records
nb_records = int(max(max(training_set[:,0]), max(test_set[:,0])))
nb_attr = len(training_set[0])

# Converting the data into Torch tensors
training_set = torch.FloatTensor(training_set)
test_set = torch.FloatTensor(test_set)
predict_set = torch.FloatTensor(predict_set)

# Creating the architecture of the Neural Network
class SAE(nn.Module):
    def __init__(self, ):
        super(SAE, self).__init__()
        self.fc1 = nn.Linear(nb_attr, 20)
        self.fc2 = nn.Linear(20, 10)
        self.fc3 = nn.Linear(10, 20)
        self.fc4 = nn.Linear(20, nb_attr)
        self.activation = nn.Sigmoid()
    def forward(self, x):
        x = self.activation(self.fc1(x))
        x = self.activation(self.fc2(x))
        x = self.activation(self.fc3(x))
        x = self.fc4(x)
        return x

def runNetwork(input_train, input_test):
	sae = SAE()
	criterion = nn.MSELoss()
	optimizer = optim.RMSprop(sae.parameters(), lr = 0.01, weight_decay = 0.5)

	# Training the SAE
	nb_epoch = 200
	for epoch in range(1, nb_epoch + 1):
		train_loss = 0
		s = 0.
		for id_record in range(nb_records):
			input = Variable(training_set[id_record]).unsqueeze(0)
			target = input.clone()
			if torch.sum(target.data > 0) > 0:
				output = sae(input)
				target.require_grad = False
				output[target == 0] = 0
				loss = criterion(output, target)
				mean_corrector = nb_attr/float(torch.sum(target.data > 0) + 1e-10)
				loss.backward()
				train_loss += np.sqrt(loss.data.item()*mean_corrector)
				s += 1.
				optimizer.step()
		#print('epoch: '+str(epoch)+' loss: '+str(train_loss/s))

	# Testing the SAE
	test_loss = 0
	s = 0.
	for id_record in range(nb_records):
		input = Variable(training_set[id_record]).unsqueeze(0)
		target = Variable(test_set[id_record])
		if torch.sum(target.data > 0) > 0:
			output = sae(input)
			target.require_grad = False
			output[(target == 0).unsqueeze(0)] = 0
			loss = criterion(output, target)
			mean_corrector = nb_attr/float(torch.sum(target.data > 0) + 1e-10)
			test_loss += np.sqrt(loss.data.item()*mean_corrector)
			s += 1.
	
	pred_y_input = Variable(predict_set[0]).unsqueeze(0)
	pred_y_output = sae(pred_y_input)
	
	return test_loss/s, pred_y_output.data

result = runNetwork(training_set, test_set)
print(result)
`

func RunNetwork(){
	trainMain := [][]float64{}
	train1 := []float64{2.1,2.2,2.5,2.8,3.8}
	train2 := []float64{1.2,1.8,1.6,1.9,2.9}
	train3 := []float64{3.4,3.5,3.8,3.7,4.7}

	trainMain = append(trainMain, train1)
	trainMain = append(trainMain, train2)
	trainMain = append(trainMain, train3)

	testMain := [][]float64{}
	test1 := []float64{2.1,2.2,2.5,2.8,3.8}
	test2 := []float64{1.2,1.8,1.6,1.9,2.9}
	test3 := []float64{3.4,3.5,3.8,3.7,4.7}

	testMain = append(testMain, test1)
	testMain = append(testMain, test2)
	testMain = append(testMain, test3)

	predMain := [][]float64{}
	pred1 := []float64{7.8,7.6,7.3,7.4,7.9}
	predMain = append(predMain, pred1)


	trainJson,_ := json.Marshal(trainMain)
	testJson,_ := json.Marshal(testMain)
	predJson,_ := json.Marshal(predMain)

	data := fmt.Sprintf(encoderScript, string(trainJson), string(testJson), string(predJson))
	cmd := exec.Command("python",  "-c", data)
	fmt.Println(cmd.Args)
	outStr, err := cmd.CombinedOutput()
	if err!=nil{
		//return error here?
	}
	outputs := strings.Split(string(outStr)," ")
	fmt.Println(outputs)
}