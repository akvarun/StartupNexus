"use client";
import styles from "./dashboard.module.css";
import { useEffect, useState } from "react";
import NewsFeed from "./newsfeed";
import { useForm } from "react-hook-form";
import { useSession } from "next-auth/react";
import UserProfile from "../user-profile/page";

export default function Dashboard() {

    const [isDropdownVisible, setDropdownVisible] = useState(false);
    const [isPostModalVisible, setPostModalVisible] = useState(false);
    const [selectedIndustries, setSelectedIndustries] = useState([]);
    const { register,handleSubmit, reset } = useForm();
    const { data: session, status } = useSession();
    const [posts, setPosts] = useState([]);
    const [selectedPostId, setSelectedPostId] = useState(null);
    const [searchTerm, setSearchTerm] = useState('');
    const [searchResults, setSearchResults] = useState([]);
    const [selectedUser, setSelectedUser] = useState(null); // Track selected user


  
    const industries = [
        "AI", "Healthcare/Health Tech", "Cybersecurity", "Internet of Things (IoT)", 
        "Fintech", "Clean Tech/Green Energy", "E-commerce/Retail Tech", "AgriTech", 
        "Robotics and Automation", "Online Education/Skill Development", "Personalized Nutrition/Wellness"
    ];

    const toggleDropdown = () => {
        setDropdownVisible((prev) => !prev);
    };

    const togglePostModal = () => {
        setPostModalVisible((prev) => !prev);
    };
        // Handle search
        const handleSearch = async (e) => {
            const term = e.target.value;
            setSearchTerm(term);
        };
    

        const handleKeyDown = async (e) => {
            if (e.key === 'Enter' && searchTerm.length > 2) {
                    try {
                        const response = await fetch(`http://localhost:8080/v1/user/search?name=${searchTerm}`);
    
                        if (!response.ok) {
                            throw new Error(`Failed to fetch search results: ${response.status}`);
                        }
                        const data = await response.json();
                    if (data.status === "success") {
                        setSearchResults(data.users); // Populate search results
                    } else {
                        setSearchResults([]);
                    }
                } catch (error) {
                    console.error("Error fetching search results:", error);
                }
            }
        };
      // Handle user profile selection
      const handleSelectUser = (user) => {
        setSelectedUser(user); // Set the selected user
        setSearchResults([]); // Clear search results after selecting
    };

    const toggleIndustry = (industry) => {
        setSelectedIndustries((prev) =>
            prev.includes(industry) ? prev.filter(i => i !== industry) : [...prev, industry]
        );
    };

    const toggleOptionsMenu = (postId) => {
        setSelectedPostId(selectedPostId === postId ? null : postId);
    };

    // create a post
    const onSubmit = async(data) => {
        const postData = {
            entrepreneur_email: session?.user?.email,
            startup_name: data.startupName,
            industry: selectedIndustries.join(", "),
            description: data.description,
            budget: parseFloat(data.budget),
            timeframe: data.timeframe,
        }

        try {
            const response = await fetch('http://localhost:8080/v1/startup/insert', {
                    method: "POST",
                    headers: {
                    "Content-Type": "application/json",
                    },
                    body: JSON.stringify(postData),
                });
            if (response.ok) {
                setPostModalVisible(false);
                reset();
            } else {
                const errorData = await response.json();
                alert(`Error: ${errorData.message}`);
            }
        }catch (error) {
            console.error('There has been a problem with your insert operation:', error);
            }
    };
    
    // fetch all posts
    useEffect(() => {
        const fetchPosts = async () => {
            try {
                const response = await fetch('http://localhost:8080/v1/startup/all');
                if (response.ok) {
                    const data = await response.json();
                    setPosts(data.startups || []);
                } else {
                    setPosts([]);
                }
            } catch (error) {
                console.error('Error fetching posts:', error);
                setPosts([]);
            }
        };

        fetchPosts();
    }, []);

    return (   
        <div>
            <div className={styles.horizontalbar}>
                <h1 className={styles.logoText}>StartupNexus</h1>
                <nav className={styles.navBar}>
                <input
                        type="text"
                        className={styles.searchBar}
                        placeholder="Search people and posts..."
                        value={searchTerm}
                        onChange={handleSearch}
                        onKeyDown={handleKeyDown}
                    />
                    <button className={`${styles.navItem} ${styles.active}`}>
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="none">
                        <path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z" fill="#00DC82"/>
                        </svg>
                        Home</button>
                    <button className={`${styles.navItem} ${styles.active}`}>
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="none">
                        <path d="M12 12c2.7 0 5-2.3 5-5s-2.3-5-5-5-5 2.3-5 5 2.3 5 5 5zm0 2c-3.9 0-7.2 2.1-9 5.2V22h18v-2.8c-1.8-3.1-5.1-5.2-9-5.2z" fill="#00DC82"/>
                    </svg>
                        My Network</button>
                        <button className={`${styles.navItem} ${styles.active}`}>
                         <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="#00DC82" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="icon icon-tabler icons-tabler-outline icon-tabler-message">
                            <path stroke="none" d="M0 0h20v20H0z" fill="none"/>
                            <path d="M8 9h8"/>
                            <path d="M8 13h6"/>
                            <path d="M18 4a3 3 0 0 1 3 3v8a3 3 0 0 1 -3 3h-5l-5 3v-3h-2a3 3 0 0 1 -3 -3v-8a3 3 0 0 1 3 -3h12z"/>
                        </svg>
                        Messaging</button>
                </nav>
                <div className={styles.profileSection} onClick={toggleDropdown}>
                    <div className={styles.profilePic}>
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="none">
                            <path d="M12 12c2.7 0 5-2.3 5-5s-2.3-5-5-5-5 2.3-5 5 2.3 5 5 5zm0 2c-3.9 0-7.2 2.1-9 5.2V22h18v-2.8c-1.8-3.1-5.1-5.2-9-5.2z" fill="#00DC82"/>
                        </svg>
                    </div>
                    {isDropdownVisible && (
                        <div className={styles.dropdownMenu}>
                            <button className={styles.dropdownItem}>View Profile</button>
                            <button className={styles.dropdownItem}>Sign Out</button>
                        </div>
                    )}
                </div>
            </div>
            
            {/* Search Results */}
            {searchResults.length > 0 && (
                <div className={styles.searchResults}>
                    {searchResults.map((user) => (
                        <div
                            key={user.id}
                            className={styles.searchResultItem}
                            onClick={() => handleSelectUser(user)}
                        >
                            {user.full_name} - {user.role}
                        </div>
                    ))}
                </div>
            )}

            {/* Render UserProfile or NewsFeed */}
            <div className={selectedUser ? styles.userProfileContainer : styles.mainContainer}>
                {selectedUser ? (
                    <UserProfile user={selectedUser} />
                ) : (
                    <>                
                <div className={styles.leftContainer}>
                    <NewsFeed />
                </div>
                <div className={styles.centerContainer}>
                    <div className={styles.startPostSection}>
                        <button className={styles.startPostButton} onClick={togglePostModal}>
                            Start a post
                        </button>
                        <div className={styles.displayPostSection}>
                            {posts.map(post => (
                                    <div key={post.id} className={styles.post}>
                                        <button className={styles.moreOptionsButton} onClick={() => toggleOptionsMenu(post.id)}>. . .</button>
                                        {selectedPostId === post.id && (
                                            <div className={styles.optionsMenu}>
                                                <button className={styles.optionButton}>
                                                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                    <path d="M12 20h9"/>
                                                    <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4Z"/>
                                                    </svg>
                                                Edit
                                                </button>
                                                <button className={styles.optionButton}>
                                                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                        <path d="M3 6h18"/>
                                                        <path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                                                        <path d="M10 11v6"/>
                                                        <path d="M14 11v6"/>
                                                        <path d="M5 6h14l-1 14H6Z"/>
                                                    </svg>
                                                    Delete</button>
                                            </div>
                                        )}
                                        <h4 className={styles.postStyle}>Startup Name: {post.startup_name}</h4>
                                        <p className={styles.postStyle}>Industry: {post.industry}</p>
                                        <p className={styles.postStyle}>Description: {post.description}</p>
                                        <p className={styles.postStyle}>Budget: {post.budget}</p>
                                        <p className={styles.postStyle}>Timeframe: {post.timeframe}</p>
                                        <hr className={styles.separator}></hr>
                                        <div className={styles.buttonContainer}>
                                            <button className={styles.likeButton}>Like</button>
                                            <button className={styles.commentButton}>Comment</button>
                                        </div>
                                    </div>

                            ))}
                        </div>
                    </div>
                </div>
                </>
                )}
            </div>
            {isPostModalVisible && (
                <div className={styles.modalOverlay}>
                    <div className={styles.postModal}>
                        <div className={styles.modalHeader}>
                            <h3>Create a post</h3>
                            <button className={styles.closeButton} onClick={togglePostModal}>✖</button>
                        </div>
                        <form onSubmit={handleSubmit(onSubmit)}>
                            <label className={styles.labelStyle}>
                                Startup Name: <br></br>
                                <input type="text" className={styles.formStyle} placeholder="Enter the name of your startup" {...register("startupName")}/>
                            </label>
                            <br></br><br></br>
                            <label className={styles.labelStyle}>
                                Industry: <br></br>
                                <div className={styles.industryButtons}>
                                {industries.map((industry) => (
                                    <button 
                                        key={industry}
                                        type="button" 
                                        className={selectedIndustries.includes(industry) ? styles.selected : styles.unselected}
                                        onClick={() => toggleIndustry(industry)}
                                    >
                                        {industry} {selectedIndustries.includes(industry) ? "✓" : "+"}
                                    </button>
                                ))}
                            </div>
                            </label>
                            <br></br>
                            <label className={styles.labelStyle}>
                                Description:
                                <textarea placeholder="What do you want to talk about?" className={styles.postInput} {...register("description")}></textarea>
                            </label>
                            <br></br><br></br>
                            <label className={styles.labelStyle}>
                                Budget: <br></br>
                                <input type="text" className={styles.formStyle} placeholder="Enter your estimated budget (e.g.5000)" {...register("budget")}/>
                            </label>
                            <br></br><br></br>
                            <label className={styles.labelStyle}>
                                Timeframe: <br></br>
                                <input type="text" className={styles.formStyle} placeholder="Select a timeframe for your project (6 or 8 months)" {...register("timeframe")}/>
                            </label>
                            <div className={styles.modalFooter}>
                                <button className={styles.postButton} type="submit">Post</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}