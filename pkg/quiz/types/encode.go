package types

import "encoding/json"

//EncodeData -
func EncodeData(dataType string, data interface{}) (dataWrapper *DataWrapper) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	dataWrapper = &DataWrapper{}
	switch dataType {
	case TYPEQUIZ:
		dataWrapper.DataType = TYPEQUIZ
	case TYPEUSER:
		dataWrapper.DataType = TYPEUSER
	case TYPESCORE:
		dataWrapper.DataType = TYPESCORE
	case TYPESTATUS:
		dataWrapper.DataType = TYPESTATUS
	case QUIZANSWER:
		dataWrapper.DataType = QUIZANSWER
	case QUIZRECEIPT:
		dataWrapper.DataType = QUIZRECEIPT
	default:
		//TODO -
		dataWrapper.DataType = dataType
	}
	dataWrapper.Data = b
	return
}
