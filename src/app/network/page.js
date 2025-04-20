'use client';

import { useEffect, useState } from "react";
import styles from "./network.module.css";

export default function Networks({ onBack }) {
    const [connections, setConnections] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchConnections = async () => {
            try {
                // Hardcoded dummy response for connections
                const dummyConnections = [
                    {
                        id: 1,
                        full_name: "John Doe",
                        role: "Software Engineer",
                        location: "San Francisco, CA",
                        linkedin_profile: "https://linkedin.com/in/johndoe",
                        profile_picture: "https://via.placeholder.com/100",
                    },
                    {
                        id: 2,
                        full_name: "Jane Smith",
                        role: "Product Manager",
                        location: "New York City, NY",
                        linkedin_profile: "https://linkedin.com/in/janesmith",
                        profile_picture: "https://via.placeholder.com/100",
                    },
                    {
                        id: 3,
                        full_name: "Alex Johnson",
                        role: "Data Scientist",
                        location: "Seattle, WA",
                        linkedin_profile: "https://linkedin.com/in/alexjohnson",
                        profile_picture: "https://via.placeholder.com/100",
                    },
                ];

                setConnections(dummyConnections);
            } catch (error) {
                console.error("Error fetching connections:", error);
                setError("Failed to load connections.");
            } finally {
                setLoading(false);
            }
        };

        fetchConnections();
    }, []);

    if (loading) return <div className={styles.loading}>Loading connections...</div>;
    if (error) return <div className={styles.error}>{error}</div>;

    return (
        <div className={styles.pageContainer}>
            <h1 className={styles.title}>My Network</h1>
            <div className={styles.connectionsContainer}>
                {connections.map((connection) => (
                    <div key={connection.id} className={styles.connectionCard}>
                        <div className={styles.profilePicture}>
                            <img
                                src={connection.profile_picture}
                                alt={`${connection.full_name}'s profile`}
                            />
                        </div>
                        <div className={styles.connectionInfo}>
                            <h2>{connection.full_name}</h2>
                            <p>{connection.role}</p>
                            <p>{connection.location}</p>
                            <a
                                href={connection.linkedin_profile}
                                target="_blank"
                                rel="noopener noreferrer"
                                className={styles.linkedinLink}
                            >
                                View LinkedIn
                            </a>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}