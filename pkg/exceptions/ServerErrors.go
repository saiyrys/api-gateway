package exceptions

import (
	"log"
)

type ServerErrors struct {
	serverErrors chan error
}

func NewServerErrors(bufferSize int) *ServerErrors {
	return &ServerErrors{
		serverErrors: make(chan error, bufferSize),
	}
}

func (se *ServerErrors) AddServerErrors(err error){
	if err != nil {
		select {
			case se.serverErrors <- err:

			default:
				log.Printf("Channel is overflow: %v", err)
		}
	}
}

func (se *ServerErrors) GetErrors() []error {
    var errors []error
    for {
        select {
        case err := <-se.serverErrors:
            errors = append(errors, err)
        default:
            return errors // канал пуст
        }
    }
}


func (se *ServerErrors) Close() {
	close(se.serverErrors)
}