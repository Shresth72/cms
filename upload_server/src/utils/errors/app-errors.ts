const STATUS_CODES = {
  OK: 200,
  BAD_REQUEST: 400,
  NOT_FOUND: 404,
  INTERNAL_SERVER_ERROR: 500,
};

class BaseError extends Error {
  statusCode: number;

  constructor(name: string, statusCode: number, description: string) {
    super(description);
    Object.setPrototypeOf(this, new.target.prototype);
    this.name = name;
    this.statusCode = statusCode;
    Error.captureStackTrace(this);
  }
}

// TODO: ADD MORE LATER
class NotFoundError extends BaseError {
  constructor(description = "not found") {
    super("not found", STATUS_CODES.NOT_FOUND, description);
  }
}

export { NotFoundError };
