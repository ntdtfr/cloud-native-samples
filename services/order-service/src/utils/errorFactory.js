// src/utils/errorFactory.js - Error factory utility
exports.createError = (statusCode, message, details = null) => {
  const error = new Error(message);
  error.statusCode = statusCode;
  error.details = details;
  return error;
};
