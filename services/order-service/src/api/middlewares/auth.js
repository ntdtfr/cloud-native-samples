// src/api/middlewares/auth.js - Authentication middleware
const jwt = require('jsonwebtoken');
const config = require('../../config');
const { formatResponse } = require('../../utils/responseFormatter');
const logger = require('../../utils/logger');

/**
 * Authenticate API requests using JWT.
 *
 * This middleware verifies the JWT token in the "Authorization" header
 * and attaches the user information to the request object.
 *
 * @param {Request} req - Express request object.
 * @param {Response} res - Express response object.
 * @param {NextFunction} next - Express next function.
 *
 * @throws {UnauthorizedError} If the token is invalid or missing.
 * @throws {InternalServerError} If there is an error during the token verification.
 */
exports.authenticate = (req, res, next) => {
  try {
    const authHeader = req.headers.authorization;

    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return formatResponse(res, 401, 'No authentication token provided');
    }

    const token = authHeader.split(' ')[1];

    if (!token) {
      return formatResponse(res, 401, 'Invalid authentication token format');
    }

    jwt.verify(token, config.jwt.secret, (err, decoded) => {
      if (err) {
        logger.warn('JWT verification failed:', err.message);
        return formatResponse(res, 401, 'Invalid or expired token');
      }

      // Attach user info to request object
      req.user = decoded;
      next();
    });
  } catch (error) {
    logger.error('Authentication error:', error);
    return formatResponse(res, 500, 'Authentication failed');
  }
};
