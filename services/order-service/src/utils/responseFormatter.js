// src/utils/responseFormatter.js - Response formatting utility
exports.formatResponse = (res, statusCode, message, data = null) => {
  const response = {
    success: statusCode < 400,
    message,
  };

  if (data) {
    response.data = data;
  }

  return res.status(statusCode).json(response);
};
