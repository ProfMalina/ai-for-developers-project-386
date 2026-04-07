import { describe, it, expect } from 'vitest';
import { ApiError, ApiValidationError, ApiNotFoundError, ApiConflictError } from './client';

describe('API Client Error Classes', () => {
  describe('ApiError', () => {
    it('should create an ApiError instance', () => {
      const error = new ApiError('Test error');

      expect(error).toBeInstanceOf(Error);
      expect(error.name).toBe('ApiError');
      expect(error.message).toBe('Test error');
    });

    it('should store data on ApiError', () => {
      const data = { error: 'TEST', message: 'Test error' };
      const error = new ApiError('Test error', data);

      expect(error.data).toEqual(data);
    });
  });

  describe('ApiValidationError', () => {
    it('should create an ApiValidationError instance', () => {
      const error = new ApiValidationError('Validation failed');

      expect(error).toBeInstanceOf(ApiError);
      expect(error.name).toBe('ApiValidationError');
    });
  });

  describe('ApiNotFoundError', () => {
    it('should create an ApiNotFoundError instance', () => {
      const error = new ApiNotFoundError('Not found');

      expect(error).toBeInstanceOf(ApiError);
      expect(error.name).toBe('ApiNotFoundError');
    });
  });

  describe('ApiConflictError', () => {
    it('should create an ApiConflictError instance', () => {
      const error = new ApiConflictError('Conflict');

      expect(error).toBeInstanceOf(ApiError);
      expect(error.name).toBe('ApiConflictError');
    });
  });
});
