// src/api/validators/orderValidator.js - Request validation
const Joi = require('joi');
const { formatResponse } = require('../../utils/responseFormatter');

// Schema for validating order items
const orderItemSchema = Joi.object({
  productId: Joi.string().required(),
  quantity: Joi.number().integer().min(1).required(),
  price: Joi.number().min(0).required(),
  name: Joi.string().required(),
});

// Schema for validating shipping address
const shippingAddressSchema = Joi.object({
  street: Joi.string().required(),
  city: Joi.string().required(),
  state: Joi.string().required(),
  country: Joi.string().required(),
  zipCode: Joi.string().required(),
});

// Schema for creating a new order
const createOrderSchema = Joi.object({
  items: Joi.array().items(orderItemSchema).min(1).required(),
  shippingAddress: shippingAddressSchema.required(),
  paymentMethod: Joi.string().valid('CREDIT_CARD', 'DEBIT_CARD', 'PAYPAL', 'BANK_TRANSFER').required(),
});

// Schema for updating order status
const updateOrderStatusSchema = Joi.object({
  status: Joi.string().valid('PROCESSING', 'SHIPPED', 'DELIVERED', 'CANCELLED').required(),
});

// Middleware for validating order creation requests
exports.validateCreateOrder = (req, res, next) => {
  const { error } = createOrderSchema.validate(req.body, {
    abortEarly: false,
    stripUnknown: true,
  });

  if (error) {
    const errorMessage = error.details.map(detail => detail.message).join(', ');
    return formatResponse(res, 400, errorMessage);
  }

  next();
};

// Middleware for validating order status updates
exports.validateUpdateOrderStatus = (req, res, next) => {
  const { error } = updateOrderStatusSchema.validate(req.body, {
    abortEarly: false,
    stripUnknown: true,
  });

  if (error) {
    const errorMessage = error.details.map(detail => detail.message).join(', ');
    return formatResponse(res, 400, errorMessage);
  }

  next();
};
