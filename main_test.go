package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestReadCSV(t *testing.T) {
	tests := []struct {
		name       string
		filePath   string
		want       [][]string
		wantHeader bool
		wantErr    bool
	}{
		{
			name:     "Valid CSV File",
			filePath: "testdata/valid.csv",
			want: [][]string{
				{"header1", "header2", "header3"},
				{"value1", "value2", "value3"},
				{"value4", "value5", "value6"},
			},
			wantHeader: true,
			wantErr:    false,
		},
		{
			name:       "Empty CSV File",
			filePath:   "testdata/empty.csv",
			want:       [][]string{},
			wantHeader: false,
			wantErr:    false,
		},
		{
			name:     "Non-existent CSV File",
			filePath: "testdata/nonexistent.csv",
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readCSV(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("readCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readCSV() got = %v, want %v", got, tt.want)
			}

			// Check if the header is skipped or not
			if tt.wantHeader {
				if len(got) == 0 || !reflect.DeepEqual(got[0], tt.want[0]) {
					t.Errorf("readCSV() header not skipped, got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestCreateTransactions(t *testing.T) {
	type args struct {
		records [][]string
	}
	tests := []struct {
		name string
		args args
		want map[string]*Transaction
	}{
		{
			name: "Empty Records",
			args: args{
				records: [][]string{},
			},
			want: map[string]*Transaction{},
		},
		{
			name: "Single Transaction",
			args: args{
				records: [][]string{
					{"tx1", "100", "200", ""},
				},
			},
			want: map[string]*Transaction{
				"tx1": {
					TxID:        "tx1",
					Fee:         100,
					Weight:      200,
					ParentTxIDs: []string{},
					TotalFee:    100,
					Children:    []*Transaction{},
				},
			},
		},
		{
			name: "Multiple Transactions",
			args: args{
				records: [][]string{
					{"tx1", "100", "200", ""},
					{"tx2", "150", "250", "tx1"},
					{"tx3", "80", "180", "tx1"},
				},
			},
			want: map[string]*Transaction{
				"tx1": {
					TxID:        "tx1",
					Fee:         100,
					Weight:      200,
					ParentTxIDs: []string{},
					TotalFee:    100,
					Children:    []*Transaction{},
				},
				"tx2": {
					TxID:        "tx2",
					Fee:         150,
					Weight:      250,
					ParentTxIDs: []string{"tx1"},
					TotalFee:    250,
					Children:    []*Transaction{},
				},
				"tx3": {
					TxID:        "tx3",
					Fee:         80,
					Weight:      180,
					ParentTxIDs: []string{"tx1"},
					TotalFee:    180,
					Children:    []*Transaction{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createTransactions(tt.args.records)
			if len(got) != len(tt.want) {
				t.Errorf("len(createTransactions()) = %v, want %v", len(got), len(tt.want))
				return
			}
			for txID, wantTx := range tt.want {
				gotTx, exists := got[txID]
				if !exists {
					t.Errorf("Transaction with txID %v not found", txID)
					return
				}
				if !equalTransactions(gotTx, wantTx) {
					t.Errorf("createTransactions() = %v, want %v", gotTx, wantTx)
				}
			}
		})
	}
}

func TestSelectTransactions(t *testing.T) {
	type args struct {
		transaction          *Transaction
		transactions         map[string]*Transaction
		selectedTransactions *[]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Select_Transactions",
			args: args{
				transaction: &Transaction{TxID: "tx1"},
				transactions: map[string]*Transaction{
					"tx1": &Transaction{TxID: "tx1", Fee: 100, TotalFee: 100},
					"tx2": &Transaction{TxID: "tx2", Fee: 150, TotalFee: 250},
					"tx3": &Transaction{TxID: "tx3", Fee: 80, TotalFee: 180},
				},
				selectedTransactions: &[]string{},
			},
			want: []string{"tx1", "tx2", "tx3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selectTransactions(tt.args.transaction, tt.args.transactions, tt.args.selectedTransactions)
			got := make([]string, len(*tt.args.selectedTransactions))
			for i, tx := range *tt.args.selectedTransactions {
				got[i] = tx
			}
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectTransactions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortAndPrintTransactions(t *testing.T) {
	type args struct {
		transactions map[string]*Transaction
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Sort_And_Print_Transactions",
			args: args{
				transactions: map[string]*Transaction{
					"tx1": &Transaction{TxID: "tx1", Fee: 100, TotalFee: 100},
					"tx2": &Transaction{TxID: "tx2", Fee: 150, TotalFee: 250},
					"tx3": &Transaction{TxID: "tx3", Fee: 80, TotalFee: 180},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortAndPrintTransactions(tt.args.transactions)
		})
	}
}

func TestSortTransactions(t *testing.T) {
	type args struct {
		transaction        *Transaction
		transactions       map[string]*Transaction
		sortedTransactions *[]*Transaction
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Sort_Transactions",
			args: args{
				transaction: &Transaction{
					TxID:     "tx1",
					TotalFee: 100,
					Children: []*Transaction{
						{TxID: "tx2", TotalFee: 150},
						{TxID: "tx3", TotalFee: 80},
					},
				},
				transactions: map[string]*Transaction{
					"tx1": &Transaction{TxID: "tx1", TotalFee: 100},
					"tx2": &Transaction{TxID: "tx2", TotalFee: 150},
					"tx3": &Transaction{TxID: "tx3", TotalFee: 80},
				},
				sortedTransactions: &[]*Transaction{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortTransactions(tt.args.transaction, tt.args.transactions, tt.args.sortedTransactions)
		})
	}
}

// helper function -  equalTransactions checks if two transactions are equal
func equalTransactions(a, b *Transaction) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.TxID == b.TxID &&
		a.Fee == b.Fee &&
		a.Weight == b.Weight &&
		reflect.DeepEqual(a.ParentTxIDs, b.ParentTxIDs) &&
		a.IsSelected == b.IsSelected &&
		reflect.DeepEqual(a.Children, b.Children)
}
