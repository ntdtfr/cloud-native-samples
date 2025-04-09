// src/api/middlewares/errorHandler.js - Error handling middleware
const logger = require('../../utils/logger');
const { formatResponse } = require('../../utils/responseFormatter');

/**
 * Error handling middleware for Express applications.
 *
 * This middleware captures errors thrown in the application and sends
 * a formatted response to the client. It logs the error details using
 * the appropriate log level based on the status code. For 5xx errors,
 * it logs at the error level, and for 4xx errors, it logs at the warn level.
 * The stack trace and additional error details are included in the response
 * only in non-production environments.
 *
 * @param {Error} err - The error object, which may contain a statusCode and message.
 * @param {Request} req - Express request object.
 * @param {Response} res - Express response object.
 * @param {NextFunction} next - Express next function.
 */

exports.errorHandler = (err, req, res, next) => {
  const statusCode = err.statusCode || 500;
  const message = err.message || 'Internal Server Error';

  // Log error details
  if (statusCode >= 500) {
    logger.error(`[${statusCode}] ${message}`, {
      path: req.path,
      method: req.method,
      error: err.stack
    });
  } else {
    logger.warn(`[${statusCode}] ${message}`, {
      path: req.path,
      method: req.method
    });
  }

  // Don't expose stack traces in production
  const errorDetails = process.env.NODE_ENV === 'production'
    ? undefined
    : { stack: err.stack, details: err.details };

  return formatResponse(res, statusCode, message, errorDetails);
};
