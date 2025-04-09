// src/api/middlewares/validation.js - General validation middleware
const { formatResponse } = require('../../utils/responseFormatter');

// Middleware for validating request parameters
exports.validateParams = (schema) => {
  return (req, res, next) => {
    const { error } = schema.validate(req.params, {
      abortEarly: false,
      stripUnknown: true,
    });

    if (error) {
      const errorMessage = error.details.map(detail => detail.message).join(', ');
      return formatResponse(res, 400, errorMessage);
    }

    next();
  };
};

// Middleware for validating request query parameters
exports.validateQuery = (schema) => {
  return (req, res, next) => {
    const { error } = schema.validate(req.query, {
      abortEarly: false,
      stripUnknown: true,
    });

    if (error) {
      const errorMessage = error.details.map(detail => detail.message).join(', ');
      return formatResponse(res, 400, errorMessage);
    }

    next();
  };
};
