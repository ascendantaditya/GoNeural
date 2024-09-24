package validation

import (
	"math/rand"
	"time"

	// third part import
	log "github.com/sirupsen/logrus"

	// internal import
	mn "github.com/ascendantaditya/goneural/model/neural"
	mu "github.com/ascendantaditya/goneural/util"
	//"fmt"
)

// TrainTestPatternsSplit split an array of patterns in training and testing.
// if shuffle is 0 the function takes the first percentage items as train and the other as test
// otherwise the patterns array is shuffled before partitioning
func TrainTestPatternsSplit(patterns []mn.Pattern, percentage float64, shuffle int) (train []mn.Pattern, test []mn.Pattern) {

	// create splitting pivot
	var splitPivot int = int(float64(len(patterns)) * percentage)
	train = make([]mn.Pattern, splitPivot)
	test = make([]mn.Pattern, len(patterns)-splitPivot)

	// if mixed mode, split with shuffling
	if shuffle == 1 {
		// create random indexes permutation
		rand.Seed(time.Now().UTC().UnixNano())
		perm := rand.Perm(len(patterns))

		// copy training data
		for i := 0; i < splitPivot; i++ {
			train[i] = patterns[perm[i]]
		}
		// copy test data
		for i := 0; i < len(patterns)-splitPivot; i++ {
			test[i] = patterns[perm[i]]
		}

	} else {
		// else, split without shuffle
		train = patterns[:splitPivot]
		test = patterns[splitPivot:]
	}

	log.WithFields(log.Fields{
		"level":     "info",
		"msg":       "splitting completed",
		"trainSet":  len(train),
		"testSet: ": len(test),
	}).Info("Complete splitting train/test set.")

	return train, test
}

// TrainTestPatternsSplit split an array of patterns in training and testing.
// if shuffle is 0 the function takes the first percentage items as train and the other as test
// otherwise the patterns array is shuffled before partitioning
func TrainTestPatternSplit(patterns []mn.Pattern, percentage float64, shuffle int) (train []mn.Pattern, test []mn.Pattern) {

	// create splitting pivot
	var splitPivot int = int(float64(len(patterns)) * percentage)
	train = make([]mn.Pattern, splitPivot)
	test = make([]mn.Pattern, len(patterns)-splitPivot)

	// if mixed mode, split with shuffling
	if shuffle == 1 {
		// create random indexes permutation
		rand.Seed(time.Now().UTC().UnixNano())
		perm := rand.Perm(len(patterns))

		// copy training data
		for i := 0; i < splitPivot; i++ {
			train[i] = patterns[perm[i]]
		}
		// copy test data
		for i := 0; i < len(patterns)-splitPivot; i++ {
			test[i] = patterns[perm[i]]
		}

	} else {
		// else, split without shuffle
		train = patterns[:splitPivot]
		test = patterns[splitPivot:]
	}

	log.WithFields(log.Fields{
		"level":     "info",
		"msg":       "splitting completed",
		"trainSet":  len(train),
		"testSet: ": len(test),
	}).Info("Complete splitting train/test set.")

	return train, test
}

// KFoldPatternsSplit split an array of patterns in k subsets.
// if shuffle is 0 the function partitions the items maintaining the order
// otherwise the patterns array is shuffled before partitioning
func KFoldPatternsSplit(patterns []mn.Pattern, k int, shuffle int) [][]mn.Pattern {

	// get the size of each fold
	var size = int(len(patterns) / k)
	var freeElements = int(len(patterns) % k)

	folds := make([][]mn.Pattern, k)

	var perm []int
	// if mixed mode, split with shuffling
	if shuffle == 1 {
		// create random indexes permutation
		rand.Seed(time.Now().UTC().UnixNano())
		perm = rand.Perm(len(patterns))
	}

	// start splitting
	currSize := 0
	foldStart := 0
	curr := 0
	for f := 0; f < k; f++ {
		curr = foldStart
		currSize = size
		if f < freeElements {
			// add another
			currSize++
		}

		// create array
		folds[f] = make([]mn.Pattern, currSize)

		// copy elements

		for i := 0; i < currSize; i++ {
			if shuffle == 1 {
				folds[f][i] = patterns[perm[curr]]
			} else {
				folds[f][i] = patterns[curr]
			}
			curr++
		}

		foldStart = curr

	}

	log.WithFields(log.Fields{
		"level":              "info",
		"msg":                "splitting completed",
		"numberOfFolds":      k,
		"meanFoldSize: ":     size,
		"consideredElements": (size * k) + freeElements,
	}).Info("Complete folds splitting.")

	return folds
}

// RandomSubsamplingValidation perform evaluation on neuron algorithm.
// It returns scores reached for each fold iteration.
func RandomSubsamplingValidation(neuron *mn.NeuronUnit, patterns []mn.Pattern, percentage float64, epochs int, folds int, shuffle int) []float64 {

	// results and predictions vars init
	var scores, actual, predicted []float64
	var train, test []mn.Pattern

	scores = make([]float64, folds)

	for t := 0; t < folds; t++ {
		// split the dataset with shuffling
		train, test = TrainTestPatternsSplit(patterns, percentage, shuffle)

		// train neuron with set of patterns, for specified number of epochs
		mn.TrainNeuron(neuron, train, epochs, 1)

		// compute predictions for each pattern in testing set
		for _, pattern := range test {
			actual = append(actual, pattern.SingleExpectation)
			predicted = append(predicted, mn.Predict(neuron, &pattern))
		}

		// compute score
		_, percentageCorrect := mn.Accuracy(actual, predicted)
		scores[t] = percentageCorrect

		log.WithFields(log.Fields{
			"level":             "info",
			"place":             "validation",
			"method":            "RandomSubsamplingValidation",
			"foldNumber":        t,
			"trainSetLen":       len(train),
			"testSetLen":        len(test),
			"percentageCorrect": percentageCorrect,
		}).Info("Evaluation completed for current fold.")
	}

	// compute average score
	acc := 0.0
	for i := 0; i < len(scores); i++ {
		acc += scores[i]
	}

	mean := acc / float64(len(scores))

	log.WithFields(log.Fields{
		"level":       "info",
		"place":       "validation",
		"method":      "RandomSubsamplingValidation",
		"folds":       folds,
		"trainSetLen": len(train),
		"testSetLen":  len(test),
		"meanScore":   mean,
	}).Info("Evaluation completed for all folds.")

	return scores
}

// RandomSubsamplingValidation perform evaluation on neuron algorithm.
// It returns scores reached for each fold iteration.
func KFoldValidation(neuron *mn.NeuronUnit, patterns []mn.Pattern, epochs int, k int, shuffle int) []float64 {

	// results and predictions vars init
	var scores, actual, predicted []float64
	var train, test []mn.Pattern

	scores = make([]float64, k)

	// split the dataset with shuffling
	folds := KFoldPatternsSplit(patterns, k, shuffle)

	// the t-th fold is used as test
	for t := 0; t < k; t++ {
		// prepare train
		train = nil
		for i := 0; i < k; i++ {
			if i != t {
				train = append(train, folds[i]...)
			}
		}
		test = folds[t]

		// train neuron with set of patterns, for specified number of epochs
		mn.TrainNeuron(neuron, train, epochs, 1)

		// compute predictions for each pattern in testing set
		for _, pattern := range test {
			actual = append(actual, pattern.SingleExpectation)
			predicted = append(predicted, mn.Predict(neuron, &pattern))
		}

		// compute score
		_, percentageCorrect := mn.Accuracy(actual, predicted)
		scores[t] = percentageCorrect

		log.WithFields(log.Fields{
			"level":             "info",
			"place":             "validation",
			"method":            "KFoldValidation",
			"foldNumber":        t,
			"trainSetLen":       len(train),
			"testSetLen":        len(test),
			"percentageCorrect": percentageCorrect,
		}).Info("Evaluation completed for current fold.")
	}

	// compute average score
	acc := 0.0
	for i := 0; i < len(scores); i++ {
		acc += scores[i]
	}

	mean := acc / float64(len(scores))

	log.WithFields(log.Fields{
		"level":       "info",
		"place":       "validation",
		"method":      "KFoldValidation",
		"folds":       k,
		"trainSetLen": len(train),
		"testSetLen":  len(test),
		"meanScore":   mean,
	}).Info("Evaluation completed for all folds.")

	return scores

}

// It returns scores reached for each fold iteration.
func MLPRandomSubsamplingValidation(mlp *mn.MultiLayerNetwork, patterns []mn.Pattern, percentage float64, epochs int, folds int, shuffle int, mapped []string) []float64 {

	// results and predictions vars init
	var scores, actual, predicted []float64
	var train, test []mn.Pattern

	scores = make([]float64, folds)

	for t := 0; t < folds; t++ {
		// split the dataset with shuffling
		train, test = TrainTestPatternsSplit(patterns, percentage, shuffle)

		// train mlp with set of patterns, for specified number of epochs
		mn.MLPTrain(mlp, patterns, mapped, epochs)

		// compute predictions for each pattern in testing set
		for _, pattern := range test {
			// get actual
			actual = append(actual, pattern.SingleExpectation)
			// get output from network
			o_out := mn.Execute(mlp, &pattern)
			// get index of max output
			_, indexMaxOut := mu.MaxInSlice(o_out)
			// add to predicted values
			predicted = append(predicted, float64(indexMaxOut))
		}

		// compute score
		_, percentageCorrect := mn.Accuracy(actual, predicted)
		scores[t] = percentageCorrect

		log.WithFields(log.Fields{
			"level":             "info",
			"place":             "validation",
			"method":            "MLPRandomSubsamplingValidation",
			"foldNumber":        t,
			"trainSetLen":       len(train),
			"testSetLen":        len(test),
			"percentageCorrect": percentageCorrect,
		}).Info("Evaluation completed for current fold.")
	}

	// compute average score
	acc := 0.0
	for i := 0; i < len(scores); i++ {
		acc += scores[i]
	}

	mean := acc / float64(len(scores))

	log.WithFields(log.Fields{
		"level":       "info",
		"place":       "validation",
		"method":      "MLPRandomSubsamplingValidation",
		"folds":       folds,
		"trainSetLen": len(train),
		"testSetLen":  len(test),
		"meanScore":   mean,
	}).Info("Evaluation completed for all folds.")

	return scores
}

// RandomSubsamplingValidation perform evaluation on neuron algorithm.
// It returns scores reached for each fold iteration.
func MLPKFoldValidation(mlp *mn.MultiLayerNetwork, patterns []mn.Pattern, epochs int, k int, shuffle int, mapped []string) []float64 {

	// results and predictions vars init
	var scores, actual, predicted []float64
	var train, test []mn.Pattern

	scores = make([]float64, k)

	// split the dataset with shuffling
	folds := KFoldPatternsSplit(patterns, k, shuffle)

	// the t-th fold is used as test
	for t := 0; t < k; t++ {
		// prepare train
		train = nil
		for i := 0; i < k; i++ {
			if i != t {
				train = append(train, folds[i]...)
			}
		}
		test = folds[t]

		// train mlp with set of patterns, for specified number of epochs
		mn.MLPTrain(mlp, patterns, mapped, epochs)

		// compute predictions for each pattern in testing set
		for _, pattern := range test {
			// get actual
			actual = append(actual, pattern.SingleExpectation)
			// get output from network
			o_out := mn.Execute(mlp, &pattern)
			// get index of max output
			_, indexMaxOut := mu.MaxInSlice(o_out)
			// add to predicted values
			predicted = append(predicted, float64(indexMaxOut))
		}

		// compute score
		_, percentageCorrect := mn.Accuracy(actual, predicted)
		scores[t] = percentageCorrect

		log.WithFields(log.Fields{
			"level":             "info",
			"place":             "validation",
			"method":            "MLPKFoldValidation",
			"foldNumber":        t,
			"trainSetLen":       len(train),
			"testSetLen":        len(test),
			"percentageCorrect": percentageCorrect,
		}).Info("Evaluation completed for current fold.")
	}

	// compute average score
	acc := 0.0
	for i := 0; i < len(scores); i++ {
		acc += scores[i]
	}

	mean := acc / float64(len(scores))

	log.WithFields(log.Fields{
		"level":       "info",
		"place":       "validation",
		"method":      "MLPKFoldValidation",
		"folds":       k,
		"trainSetLen": len(train),
		"testSetLen":  len(test),
		"meanScore":   mean,
	}).Info("Evaluation completed for all folds.")

	return scores

}

// RNNValidation perform evaluation on neuron algorithm.
func RNNValidation(mlp *mn.MultiLayerNetwork, patterns []mn.Pattern, epochs int, shuffle int) (float64, []float64) {

	// results and predictions vars init
	var scores []float64
	scores = make([]float64, len(patterns))

	// train mlp with set of patterns, for specified number of epochs
	mn.ElmanTrain(mlp, patterns, epochs)
	p_cor := 0.0

	// compute predictions for each pattern in testing set
	for p_i, pattern := range patterns {
		// get output from network
		o_out := mn.Execute(mlp, &pattern, 1)
		for o_out_i, o_out_v := range o_out {
			o_out[o_out_i] = mu.Round(o_out_v, .5, 0)
		}
		log.WithFields(log.Fields{
			"a_p_b": pattern.Features,
			"rea_c": pattern.MultipleExpectation,
			"pre_c": o_out,
		}).Debug()

		// add to predicted values
		_, p_cor = mn.Accuracy(pattern.MultipleExpectation, o_out)
		// compute score
		scores[p_i] = p_cor
	}

	// compute average score
	acc := 0.0
	for i := 0; i < len(scores); i++ {
		acc += scores[i]
	}

	mean := acc / float64(len(scores))

	log.WithFields(log.Fields{
		"level":       "info",
		"place":       "validation",
		"method":      "RNNValidation",
		"trainSetLen": len(patterns),
		"testSetLen":  len(patterns),
		"meanScore":   mean,
	}).Info("Evaluation completed for all patterns.")

	return mean, scores

}
