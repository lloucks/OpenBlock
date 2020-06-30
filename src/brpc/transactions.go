//RPCs in here will be oriented around recieving transactions to put into blocks


package brpc

import (
    "structures"
)



//Basic rpc structures to communicate transactions
type Args struct{
    Transaction structures.Transaction

}


type Reply struct{
    valid bool

}
