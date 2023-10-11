package constants

import helpers "kanban-board/helpers/secrets"

var SECRET_JWT_KEY, _ = helpers.LoadSecrets(".env", "SECRETJTW")