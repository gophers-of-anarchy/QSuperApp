package messages

const InternalError = "Internal server error"
const InvalidRequestBody = "Invalid request body"
const InvalidCardNumber = "Invalid card number"
const InvalidTransferAmount = "Invalid transfer amount"
const InsufficientBalance = "Insufficient balance"
const FailedToUpdateSender = "Failed to update sender's card balance"
const FailedToUpdateReceiver = "Failed to update receiver's card balance"
const FailedToCreateTransaction = "Failed to create transaction record"
const FailedToCreateTransactionFee = "Failed to create transaction fee record"
const FailedToCommitTransaction = "Failed to commit transaction"
const FailedToSendSenderSMS = "Failed to send sender's sms"
const FailedToSendReceiverSMS = "Failed to send receiver's sms"
const FailedToEncodeResponse = "Failed to encode response"
const TooManyRequests = "Too many requests"

const AccountNotFound = "Account not found"
const FailedToCreateCard = "Failed to create card"
const FailedToCreateAccount = "Failed to create account"
const FailedToCreateOrder = "Failed to create order"

const FailedPasswordHashGeneration = "Failed to generate password hash"
const FailedToCreateUser = "Failed to create user"
const UsernameOrPasswordIncorrect = "Username or password is incorrect"
const FailedToCreateToken = "Failed to create token"
const Unauthorized = "Unauthorized"
const InvalidToken = "Invalid token"
