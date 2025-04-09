// tests/unit/orderService.test.js - Unit test for order service
const mongoose = require('mongoose');
const { MongoMemoryServer } = require('mongodb-memory-server');
const Order = require('../../src/models/Order');
const orderService = require('../../src/services/orderService');

let mongoServer;

beforeAll(async () => {
  mongoServer = await MongoMemoryServer.create();
  await mongoose.connect(mongoServer.getUri(), {
    useNewUrlParser: true,
    useUnifiedTopology: true,
  });
});

afterAll(async () => {
  await mongoose.disconnect();
  await mongoServer.stop();
});

beforeEach(async () => {
  await Order.deleteMany({});
});

describe('Order Service', () => {
  const customerId = 'customer123';

  const sampleOrder = {
    customerId,
    items: [
      {
        productId: 'product123',
        quantity: 2,
        price: 29.99,
        name: 'Test Product',
      },
    ],
    shippingAddress: {
      street: '123 Test St',
      city: 'Test City',
      state: 'Test State',
      country: 'Test Country',
      zipCode: '12345',
    },
    paymentMethod: 'CREDIT_CARD',
  };

  describe('createOrder', () => {
    it('should create a new order', async () => {
      const result = await orderService.createOrder(sampleOrder);

      expect(result).toBeDefined();
      expect(result.customerId).toBe(customerId);
      expect(result.items.length).toBe(1);
      expect(result.status).toBe('PENDING');
      expect(result.paymentStatus).toBe('PENDING');
      expect(result.totalAmount).toBe(59.98); // 29.99 * 2
    });

    it('should throw an error if required fields are missing', async () => {
      const invalidOrder = { ...sampleOrder };
      delete invalidOrder.items;

      await expect(orderService.createOrder(invalidOrder)).rejects.toThrow();
    });
  });

  describe('getOrderById', () => {
    it('should return the order if it exists and belongs to the customer', async () => {
      const newOrder = new Order(sampleOrder);
      await newOrder.save();

      const result = await orderService.getOrderById(newOrder._id, customerId);

      expect(result).toBeDefined();
      expect(result._id.toString()).toBe(newOrder._id.toString());
    });

    it('should return null if the order does not belong to the customer', async () => {
      const newOrder = new Order(sampleOrder);
      await newOrder.save();

      const result = await orderService.getOrderById(newOrder._id, 'wrong-customer');

      expect(result).toBeNull();
    });
  });

  describe('updateOrderStatus', () => {
    it('should update the order status', async () => {
      const newOrder = new Order(sampleOrder);
      await newOrder.save();

      const result = await orderService.updateOrderStatus(
        newOrder._id,
        'SHIPPED',
        customerId
      );

      expect(result).toBeDefined();
      expect(result.status).toBe('SHIPPED');
    });

    it('should throw an error when updating a cancelled order', async () => {
      const newOrder = new Order({
        ...sampleOrder,
        status: 'CANCELLED',
      });
      await newOrder.save();

      await expect(
        orderService.updateOrderStatus(newOrder._id, 'SHIPPED', customerId)
      ).rejects.toThrow('Cannot update a cancelled order');
    });
  });

  describe('cancelOrder', () => {
    it('should cancel an order and update payment status', async () => {
      const newOrder = new Order(sampleOrder);
      await newOrder.save();

      const result = await orderService.cancelOrder(newOrder._id, customerId);

      expect(result).toBeDefined();
      expect(result.status).toBe('CANCELLED');
      expect(result.paymentStatus).toBe('REFUNDED');
    });

    it('should throw an error when cancelling a delivered order', async () => {
      const newOrder = new Order({
        ...sampleOrder,
        status: 'DELIVERED',
      });
      await newOrder.save();

      await expect(
        orderService.cancelOrder(newOrder._id, customerId)
      ).rejects.toThrow('Cannot cancel a delivered order');
    });
  });
});
