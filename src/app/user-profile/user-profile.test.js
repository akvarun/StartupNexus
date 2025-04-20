import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import UserProfile from "./page";
import { useSession } from "next-auth/react";
import '@testing-library/jest-dom/extend-expect';

// Mocking next-auth useSession hook
jest.mock("next-auth/react", () => ({
  useSession: jest.fn(),
}));

// Mocking fetch function
global.fetch = jest.fn();

describe("UserProfile", () => {
  const user = {
    id: "1",
    full_name: "John Doe",
    email: "johndoe@example.com",
  };

  beforeEach(() => {
    useSession.mockReturnValue({ data: { user: { email: "sessionUser@example.com" } } });
    fetch.mockClear();
  });

  it("renders loading state initially", () => {
    render(<UserProfile user={user} />);
    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("renders user data after successful fetch", async () => {
    const mockResponse = {
      ok: true,
      json: async () => ({
        id: "1",
        name: "John Doe",
        role: "Full Stack Developer",
        location: "San Francisco, CA",
        about: "Passionate about building scalable web applications and microservices.",
        activity: [{ id: 1, content: "Exploring the power of Kubernetes!" }],
        experience: [{ id: 1, company: "Google", role: "Software Engineer", duration: "2 yrs" }],
        profilePicture: "https://example.com/profile-picture.jpg",
      }),
    };

    fetch.mockResolvedValue(mockResponse);

    render(<UserProfile user={user} />);

    await waitFor(() => expect(fetch).toHaveBeenCalled());

    // Check for user data rendering
    expect(screen.getByText("John Doe")).toBeInTheDocument();
    expect(screen.getByText("Full Stack Developer")).toBeInTheDocument();
    expect(screen.getByText("San Francisco, CA")).toBeInTheDocument();
    expect(screen.getByText("Passionate about building scalable web applications and microservices.")).toBeInTheDocument();
  });

  it("displays error message when data fetch fails", async () => {
    const mockResponse = {
      ok: false,
      json: async () => ({ message: "Failed to load user data" }),
    };

    fetch.mockResolvedValue(mockResponse);

    render(<UserProfile user={user} />);

    await waitFor(() => expect(fetch).toHaveBeenCalled());

    // Check for error rendering
    expect(screen.getByText("Error: Failed to load user data")).toBeInTheDocument();
  });

  it("displays 'User not found' when no user data is available", async () => {
    const mockResponse = {
      ok: true,
      json: async () => null,
    };

    fetch.mockResolvedValue(mockResponse);

    render(<UserProfile user={user} />);

    await waitFor(() => expect(fetch).toHaveBeenCalled());

    // Check for no data message
    expect(screen.getByText("User not found")).toBeInTheDocument();
  });

  it("shows info box message when connect button is clicked and user is not logged in", async () => {
    useSession.mockReturnValue({ data: null });

    render(<UserProfile user={user} />);

    const connectButton = screen.getByText("Connect");
    fireEvent.click(connectButton);

    // Check for info box message
    expect(screen.getByText("You must be logged in to send a connection request.")).toBeInTheDocument();
  });

  it("sends connection request successfully", async () => {
    const mockResponse = {
      ok: true,
      json: async () => ({ message: "Success" }),
    };

    fetch.mockResolvedValue(mockResponse);

    render(<UserProfile user={user} />);

    const connectButton = screen.getByText("Connect");
    fireEvent.click(connectButton);

    await waitFor(() => expect(fetch).toHaveBeenCalled());

    // Check for success message
    expect(screen.getByText("Connection request sent successfully to John Doe")).toBeInTheDocument();
  });

  it("shows error message if connection request fails", async () => {
    const mockResponse = {
      ok: false,
      json: async () => ({ message: "Failed to send request" }),
    };

    fetch.mockResolvedValue(mockResponse);

    render(<UserProfile user={user} />);

    const connectButton = screen.getByText("Connect");
    fireEvent.click(connectButton);

    await waitFor(() => expect(fetch).toHaveBeenCalled());

    // Check for error message
    expect(screen.getByText("Failed to send connection request. Error: Failed to send request")).toBeInTheDocument();
  });

  it("opens message dialog when message button is clicked", () => {
    const spy = jest.spyOn(window, "alert").mockImplementation(() => {});

    render(<UserProfile user={user} />);

    const messageButton = screen.getByText("Message");
    fireEvent.click(messageButton);

    expect(spy).toHaveBeenCalledWith("Opening message dialog with John Doe");

    spy.mockRestore();
  });
});
