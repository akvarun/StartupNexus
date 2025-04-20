import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect'; // for the `toBeInTheDocument` matcher
import SignupForm from './signup_form';
import { useSession } from 'next-auth/react';
import { isValidPhoneNumber } from 'libphonenumber-js';

// Mock the useSession hook
jest.mock('next-auth/react', () => ({
  useSession: jest.fn(),
}));

// Mock the libphonenumber-js phone validation function
jest.mock('libphonenumber-js', () => ({
  isValidPhoneNumber: jest.fn(),
}));

describe('SignupForm', () => {
  beforeEach(() => {
    // Mock the session
    useSession.mockReturnValue({
      data: {
        user: {
          name: 'Test User',
          email: 'testuser@example.com',
        },
      },
    });
  });

  test('renders the form with pre-filled values', () => {
    render(<SignupForm />);

    // Check if the pre-filled values are rendered correctly
    expect(screen.getByPlaceholderText('Enter your full name')).toHaveValue('Test User');
    expect(screen.getByPlaceholderText('Enter your email')).toHaveValue('testuser@example.com');
  });

  test('validates required fields', async () => {
    render(<SignupForm />);

    // Trigger the submit event
    fireEvent.submit(screen.getByText('Submit'));

    // Check for validation errors
    expect(await screen.findByText('Full Name is required')).toBeInTheDocument();
    expect(await screen.findByText('Email is required')).toBeInTheDocument();
    expect(await screen.findByText('Contact number is required')).toBeInTheDocument();
  });

  test('validates email format', async () => {
    render(<SignupForm />);

    // Input invalid email
    fireEvent.change(screen.getByPlaceholderText('Enter your email'), { target: { value: 'invalid-email' } });

    // Trigger the submit event
    fireEvent.submit(screen.getByText('Submit'));

    // Check for email validation error
    expect(await screen.findByText('Invalid email format')).toBeInTheDocument();
  });

  test('validates phone number format', async () => {
    // Mock the phone number validation
    isValidPhoneNumber.mockReturnValue(false); // Invalid phone number

    render(<SignupForm />);

    // Input an invalid phone number
    fireEvent.change(screen.getByPlaceholderText('Enter your contact number'), { target: { value: '12345' } });

    // Trigger the submit event
    fireEvent.submit(screen.getByText('Submit'));

    // Check for phone number validation error
    expect(await screen.findByText('Invalid phone number')).toBeInTheDocument();
  });

  test('submits form successfully', async () => {
    // Mock the phone number validation to be true
    isValidPhoneNumber.mockReturnValue(true); // Valid phone number

    // Mock the fetch response for the API
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      json: jest.fn().mockResolvedValue({}),
    });

    render(<SignupForm />);

    // Fill in the form with valid data
    fireEvent.change(screen.getByPlaceholderText('Enter your full name'), { target: { value: 'Test User' } });
    fireEvent.change(screen.getByPlaceholderText('Enter your email'), { target: { value: 'testuser@example.com' } });
    fireEvent.change(screen.getByPlaceholderText('Enter your contact number'), { target: { value: '1234567890' } });

    // Submit the form
    fireEvent.submit(screen.getByText('Submit'));

    // Wait for the success message
    await waitFor(() => expect(screen.getByText('User profile is created successfully. Redirecting to login...')).toBeInTheDocument());
  });

  test('shows an error when submission fails', async () => {
    // Mock the fetch response for a failed API request
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      json: jest.fn().mockResolvedValue({ message: 'Error submitting form' }),
    });

    render(<SignupForm />);

    // Fill in the form with valid data
    fireEvent.change(screen.getByPlaceholderText('Enter your full name'), { target: { value: 'Test User' } });
    fireEvent.change(screen.getByPlaceholderText('Enter your email'), { target: { value: 'testuser@example.com' } });
    fireEvent.change(screen.getByPlaceholderText('Enter your contact number'), { target: { value: '1234567890' } });

    // Submit the form
    fireEvent.submit(screen.getByText('Submit'));

    // Wait for the error message
    await waitFor(() => expect(screen.getByText('Error: Error submitting form')).toBeInTheDocument());
  });
});