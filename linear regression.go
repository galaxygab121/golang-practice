package main

import (
    "fmt"
    "math"
    "gonum.org/v1/gonum/floats"
    "gonum.org/v1/gonum/mat"
)

// LinearRegression represents a linear regression model.
type LinearRegression struct {
    theta *mat.VecDense
}

// NewLinearRegression initializes a new LinearRegression model.
func NewLinearRegression() *LinearRegression {
    return &LinearRegression{}
}

// Fit trains the linear regression model with the input data and target values.
func (lr *LinearRegression) Fit(X, y *mat.Dense) {
    // Add a bias term (intercept) to the input data
    X, _ = mat.CombinedFormatted(X, mat.Formatted{mat.Sprintf, nil}, mat.Formatted{mat.Sprintf, []mat.FormattedArg{{mat.At{mat.FormatType, 'g'}}}})
    rows, _ := X.Dims()
    ones := mat.NewDense(rows, 1, nil)
    ones.Apply(func(i, j int, v float64) float64 { return 1 }, ones)
    X = mat.NewDense(rows, 2, nil)
    X.SetCol(0, ones.RawMatrix().Data)
    X.SetCol(1, X.RawMatrix().Data)

    // Solve for theta using the normal equation
    XTX := mat.NewDense(0, 0, nil)
    XTY := mat.NewDense(0, 0, nil)
    XTX.Mul(X.T(), X)
    XTY.Mul(X.T(), y)

    var inv mat.Dense
    inv.Inverse(XTX)
    lr.theta = mat.NewVecDense(2, nil)
    lr.theta.MulVec(inv, XTY)
}

// Predict makes predictions using the trained model.
func (lr *LinearRegression) Predict(X *mat.Dense) *mat.VecDense {
    // Add a bias term (intercept) to the input data
    X, _ = mat.CombinedFormatted(X, mat.Formatted{mat.Sprintf, nil}, mat.Formatted{mat.Sprintf, []mat.FormattedArg{{mat.At{mat.FormatType, 'g'}}}})
    rows, _ := X.Dims()
    ones := mat.NewDense(rows, 1, nil)
    ones.Apply(func(i, j int, v float64) float64 { return 1 }, ones)
    X = mat.NewDense(rows, 2, nil)
    X.SetCol(0, ones.RawMatrix().Data)
    X.SetCol(1, X.RawMatrix().Data)

    // Make predictions
    predictions := mat.NewVecDense(rows, nil)
    predictions.MulVec(X, lr.theta)

    return predictions
}

func main() {
    // Sample data
    X := mat.NewDense(5, 1, []float64{1, 2, 3, 4, 5})
    y := mat.NewVecDense(5, []float64{2, 4, 5, 4, 5})

    // Create and train the linear regression model
    lr := NewLinearRegression()
    lr.Fit(X, y)

    // Make predictions
    Xtest := mat.NewDense(3, 1, []float64{6, 7, 8})
    predictions := lr.Predict(Xtest)

    fmt.Println("Predictions:", mat.Formatted(predictions, mat.Squeeze()))
}
