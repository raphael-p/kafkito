package config

const PORT_ENVAR string = "KAFKITO_PORT"
const DEFAULT_PORT string = "8083"
const MAX_QUEUES uint32 = 3
const MAX_QUEUE_LENGTH uint32 = 4
const MAX_QUEUE_NAME_BYTES uint32 = 10
const MAX_MESSAGE_HEADER_BYTES uint32 = 10
const MAX_MESSAGE_BODY_BYTES uint32 = 50
const MESSAGE_TTL int64 = 60
const MESSAGE_BATCH_SIZE uint32 = 5
