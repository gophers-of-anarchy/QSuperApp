package messages

const InternalError = "Internal server error"
const InvalidRequestBody = "Invalid request body"
const InvalidCardNumber = "Invalid card number"
const InvalidTransferAmount = "Invalid transfer amount"
const InsufficientBalance = "Insufficient balance"
const FailedToUpdateSender = "Failed to update sender's card balance"
const FailedToNotifyUser = "Failed to notify customer"
const FailedToSendEmail = "Failed to send email"
const FailedToChangeStatus = "Failed to change customer's order status"
const FailedToUpdateReceiver = "Failed to update receiver's card balance"
const FailedToCreateTransaction = "Failed to create transaction record"
const FailedToCreateTransactionFee = "Failed to create transaction fee record"
const FailedToCommitTransaction = "Failed to commit transaction"
const FailedToSendSenderSMS = "Failed to send sender's sms"
const FailedToSendReceiverSMS = "Failed to send receiver's sms"
const FailedToEncodeResponse = "Failed to encode response"
const TooManyRequests = "Too many requests"
const WrongStatus = "Wrong status for order"
const WrongOrderID = "Wrong ID for order"

const UserNotFound = "User not found"
const AccountNotFound = "Account not found"
const FailedToCreateCard = "Failed to create card"
const FailedToCreateAccount = "Failed to create account"

const FailedPasswordHashGeneration = "Failed to generate password hash"
const FailedToCreateUser = "Failed to create user"
const FailedToCreateAuthTokens = "Failed to create auth tokens"
const UsernameOrPasswordIncorrect = "Username or password is incorrect"
const FailedToCreateToken = "Failed to create token"
const Unauthorized = "Unauthorized"
const InvalidToken = "Invalid token"

const OrderNotFound = "Order not found"
const PaymentError = "Error making payment"
const PaymentNotFound = "Payment not found"
const PaymentVerificationError = "Error verifying payment"

const InvalidUserType = "Invalid user type"
const ErrUnauthorizedAccess = "User does not have access"
