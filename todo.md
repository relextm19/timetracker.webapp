1. change the sessions so they hold keyHash instead of userID and send data to frontend grouped by the keyHash in a map as in [key]data
2. create another table for tokens and store a token hash instead of plain token and update the backend code to match that
