// src/services/orderService.js - Order service
const Order = require('../models/Order');
const { createError } = require('../utils/errorFactory');
const logger = require('../utils/logger');

/**
 * Create a new order with the given items, customer ID, and payment ID.
 * @param {object} orderData - Order data
 * @param {string} orderData.items - Array of item IDs
 * @param {string} orderData.paymentId - Payment ID
 * @param {string} orderData.customerId - The ID of the customer
 * @returns {object} 201 - The created order
 * @throws {AppError} 400 - If the request body is invalid
 * @throws {AppError} 500 - If there is an error creating the order
 */
exports.createOrder = async (orderData) => {
  try {
    const order = new Order(orderData);
    await order.save();
    return order;
  } catch (error) {
    logger.error('Error creating order:', error);
    throw createError(400, 'Failed to create order', error);
  }
};

/**
 * Retrieve an order by order ID for the logged-in customer
 * @param {string} orderId - The ID of the order
 * @param {string} customerId - The ID of the customer
 * @returns {object|null} The retrieved order or null if the order is not found
 * @throws {AppError} 500 - If there is an error retrieving the order
 */
exports.getOrderById = async (orderId, customerId) => {
  try {
    return await Order.findOne({ _id: orderId, customerId });
  } catch (error) {
    logger.error(`Error retrieving order ${orderId}:`, error);
    throw createError(500, 'Failed to retrieve order', error);
  }
};

/**
 * Retrieve all orders for the customer, filtered by optional status
 * @param {string} customerId - The ID of the customer
 * @param {object} options - Query options
 * @param {number} [options.page=1] - The page number for pagination
 * @param {number} [options.limit=10] - The number of orders per page
 * @param {string} [options.status] - Optional status to filter orders
 * @returns {object} An object containing the list of orders and pagination details
 * @throws {AppError} 500 - If there is an error retrieving the orders
 */
exports.getCustomerOrders = async (customerId, options) => {
  try {
    const query = { customerId };

    if (options.status) {
      query.status = options.status;
    }

    const skip = (options.page - 1) * options.limit;

    const [orders, total] = await Promise.all([
      Order.find(query)
        .sort({ createdAt: -1 })
        .skip(skip)
        .limit(options.limit),
      Order.countDocuments(query),
    ]);

    return { orders, total };
  } catch (error) {
    logger.error(`Error retrieving orders for customer ${customerId}:`, error);
    throw createError(500, 'Failed to retrieve orders', error);
  }
};

/**
 * Update the status of an order by order ID for the logged-in customer
 * @param {string} orderId - The ID of the order
 * @param {string} status - The new status to update
 * @param {string} customerId - The ID of the customer
 * @returns {object|null} The updated order or null if the order is not found
 * @throws {AppError} 400 - If the order is cancelled or delivered and cannot be updated
 * @throws {AppError} 500 - If there is an error updating the order status
 */
exports.updateOrderStatus = async (orderId, status, customerId) => {
  try {
    const order = await Order.findOne({ _id: orderId, customerId });

    if (!order) {
      return null;
    }

    // Business logic validations
    if (order.status === 'CANCELLED') {
      throw createError(400, 'Cannot update a cancelled order');
    }

    if (order.status === 'DELIVERED' && status !== 'CANCELLED') {
      throw createError(400, 'Cannot update a delivered order');
    }

    // Update order status
    order.status = status;

    // Additional business logic based on status
    if (status === 'CANCELLED') {
      order.paymentStatus = 'REFUNDED';
    }

    await order.save();
    return order;
  } catch (error) {
    logger.error(`Error updating order ${orderId} status:`, error);
    throw createError(500, 'Failed to update order status', error);
  }
};

/**
 * Cancel an order by order ID for the logged-in customer
 * @param {string} orderId - The ID of the order
 * @param {string} customerId - The ID of the customer
 * @returns {object|null} The cancelled order or null if the order is not found
 * @throws {AppError} 400 - If the order is already cancelled or delivered
 * @throws {AppError} 500 - If there is an error cancelling the order
 */
exports.cancelOrder = async (orderId, customerId) => {
  try {
    const order = await Order.findOne({ _id: orderId, customerId });

    if (!order) {
      return null;
    }

    // Business logic validations
    if (order.status === 'CANCELLED') {
      throw createError(400, 'Order is already cancelled');
    }

    if (order.status === 'DELIVERED') {
      throw createError(400, 'Cannot cancel a delivered order');
    }

    // Update order status
    order.status = 'CANCELLED';
    order.paymentStatus = 'REFUNDED';

    await order.save();
    return order;
  } catch (error) {
    logger.error(`Error cancelling order ${orderId}:`, error);
    throw createError(500, 'Failed to cancel order', error);
  }
};
