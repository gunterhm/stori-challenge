package txnprocessor

// ITxnProcessorService is a service interface for processing incoming csv files containing account transaction data
type ITxnProcessorService interface {
	StartProcess() error
}
