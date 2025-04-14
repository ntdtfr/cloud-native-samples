// src/api/controllers/orderController.js - Order controller
const orderService = require('../../services/orderService');
const logger = require('../../utils/logger');
const { formatResponse } = require('../../utils/responseFormatter');

/**
 * POST /orders
 * @summary Create a new order for a customer
 * @description Creates a new order with the given items, customer ID, and payment ID.
 * @param {object} request.body - Order data
 * @param {string} request.body.items - Array of item IDs
 * @param {string} request.body.paymentId - Payment ID
 * @returns {object} 201 - The created order
 * @throws {AppError} 400 - If the request body is invalid
 * @throws {AppError} 401 - If the user is unauthorized
 * @throws {AppError} 500 - If there is an error creating the order
 */
exports.createOrder = async (req, res, next) => {
  try {
    const customerId = req.user.id; // Extracted from JWT token
    const body = req.body;
    const orderData = { body, customerId };

    const order = await orderService.createOrder(orderData);

    logger.info(`Order created: ${order._id} for customer: ${customerId}`);
    return formatResponse(res, 201, 'Order created successfully', { order });
  } catch (error) {
    next(error);
  }
};


/**
 * GET /orders/:id
 * @summary Retrieve an order by order ID for the logged-in customer
 * @param {string} id - The ID of the order
 * @returns {object} The retrieved order
 * @throws {AppError} If the order is not found or the user is unauthorized
 */
exports.getOrderById = async (req, res, next) => {
  try {
    const { id } = req.params;
    const customerId = req.user.id;

    const order = await orderService.getOrderById(id, customerId);

    if (!order) {
      return formatResponse(res, 404, 'Order not found');
    }

    return formatResponse(res, 200, 'Order retrieved successfully', { order });
  } catch (error) {
    next(error);
  }
};

/**
 * GET /orders
 * @summary Retrieve all orders for the logged-in customer
 * @description Fetches a paginated list of orders for the customer, filtered by optional status.
 * @param {object} request.query - Query parameters
 * @param {number} [request.query.page=1] - The page number for pagination
 * @param {number} [request.query.limit=10] - The number of orders per page
 * @param {string} [request.query.status] - Optional status to filter orders
 * @returns {object} 200 - An object containing the list of orders and pagination details
 * @throws {AppError} 500 - If there is an error retrieving the orders
 */
exports.getCustomerOrders = async (req, res, next) => {
  try {
    const customerId = req.user.id;
    const { page = 1, limit = 10, status } = req.query;

    const options = {
      page: parseInt(page, 10),
      limit: parseInt(limit, 10),
      status,
    };

    const result = await orderService.getCustomerOrders(customerId, options);

    return formatResponse(res, 200, 'Orders retrieved successfully', {
      orders: result.orders,
      pagination: {
        total: result.total,
        page: options.page,
        limit: options.limit,
        pages: Math.ceil(result.total / options.limit),
      },
    });
  } catch (error) {
    next(error);
  }
};

/**
 * PATCH /orders/:id
 * @summary Update the status of an order by order ID for the logged-in customer
 * @description Updates the status of an existing order for the customer.
 * @param {string} id - The ID of the order
 * @param {object} request.body - Order data
 * @param {string} request.body.status - The new status to update
 * @returns {object} 200 - The updated order
 * @throws {AppError} 400 - If the request body is invalid
 * @throws {AppError} 401 - If the user is unauthorized
 * @throws {AppError} 404 - If the order is not found
 * @throws {AppError} 500 - If there is an error updating the order
 */
exports.updateOrderStatus = async (req, res, next) => {
  try {
    const { id } = req.params;
    const { status } = req.body;
    const customerId = req.user.id;

    const order = await orderService.updateOrderStatus(id, status, customerId);

    if (!order) {
      return formatResponse(res, 404, 'Order not found');
    }

    logger.info(`Order ${id} status updated to ${status}`);
    return formatResponse(res, 200, 'Order status updated successfully', { order });
  } catch (error) {
    next(error);
  }
};

/**
 * DELETE /orders/:id
 * @summary Cancel an order by order ID for the logged-in customer
 * @description Cancels an existing order for the customer.
 * @param {string} id - The ID of the order
 * @returns {object} 200 - The cancelled order
 * @throws {AppError} 404 - If the order is not found
 * @throws {AppError} 500 - If there is an error cancelling the order
 */
exports.cancelOrder = async (req, res, next) => {
  try {
    const { id } = req.params;
    const customerId = req.user.id;

    const order = await orderService.cancelOrder(id, customerId);

    if (!order) {
      return formatResponse(res, 404, 'Order not found');
    }

    logger.info(`Order ${id} cancelled by customer ${customerId}`);
    return formatResponse(res, 200, 'Order cancelled successfully', { order });
  } catch (error) {
    next(error);
  }
};
