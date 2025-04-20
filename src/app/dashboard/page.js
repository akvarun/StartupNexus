"use client";
import styles from "./dashboard.module.css";
import { useEffect, useState } from "react";
import NewsFeed from "./newsfeed";
import { useForm } from "react-hook-form";
import { useSession, signOut } from "next-auth/react";
import UserProfile from "../user-profile/page";
import Networks from "../network/page";
import moment from "moment/moment";
import ViewProfile from "../view-profile/page";

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
    const [selectedUser, setSelectedUser] = useState(null); 
    const [showNetworksPage, setShowNetworksPage] = useState(false);
    const [editingPostId, setEditingPostId] = useState(null);
    const [editedPostData, setEditedPostData] = useState({});
    const industries = [
        "AI", "Healthcare/Health Tech", "Cybersecurity", "Internet of Things (IoT)", 
        "Fintech", "Clean Tech/Green Energy", "E-commerce/Retail Tech", "AgriTech", 
        "Robotics and Automation", "Online Education/Skill Development", "Personalized Nutrition/Wellness"
    ];
    const [showViewProfilePage, setShowViewProfilePage] = useState(false);
    const [likedPosts, setLikedPosts] = useState(new Set());
    const [likeCounts, setLikeCounts] = useState({});
    const [commentText, setCommentText] = useState('');
    const [comments, setComments] = useState({});
    const [openCommentPostId, setOpenCommentPostId] = useState(null);
    const [commentCounts, setCommentCounts] = useState({});
    const [visibleCommentsPostId, setVisibleCommentsPostId] = useState(null);
    const [reportModalPostId, setReportModalPostId] = useState(null);
    const [reportReason, setReportReason] = useState('');

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
        setShowNetworksPage(false);
        setShowViewProfilePage(false); // Reset view profile page when viewing another user
    };
    // Show My Network page
    const handleMyNetworkClick = () => {
        setShowNetworksPage(true); // Show My Network page
        setSelectedUser(null); // Clear selected user if any
        setShowViewProfilePage(false); // Reset view profile page
    };

    const toggleIndustry = (industry) => {
        setSelectedIndustries((prev) =>
            prev.includes(industry) ? prev.filter(i => i !== industry) : [...prev, industry]
        );
    };

    const toggleOptionsMenu = (postId) => {
        setSelectedPostId(selectedPostId === postId ? null : postId);
    };

    const handleEditClick = (post) => {
        setEditingPostId(post.id);
        setEditedPostData({
            startup_name: post.startup_name,
            industry: post.industry,
            description: post.description,
            budget: post.budget,
            timeframe: post.timeframe
        });
    };
    
    const handleEditChange = (e, field) => {
        setEditedPostData((prev) => ({
            ...prev,
            [field]: e.target.value,
        }));
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

    const handleSaveEdit = async (postId) => {
        try {
            const response = await fetch('http://localhost:8080/v1/startup/update', {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ startup_id: postId, email: session?.user?.email, ...editedPostData }),
            });
    
            if (response.ok) {
                setEditingPostId(null);
                fetchPosts(); // Refresh posts after update
            } else {
                alert("Failed to update post.");
            }
        } catch (error) {
            console.error("Error updating post:", error);
        }
    };

    // delete post
    const handleDeletePost = async (postId) => {
        try {
          const response = await fetch('http://localhost:8080/v1/startup/delete', {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              email: session?.user?.email, // Include user's email for authentication
              startup_id: postId, // Pass the startup ID to delete
            }),
          });
      
          if (response.ok) {
            // Remove the deleted post from the screen
            setPosts((prevPosts) => prevPosts.filter((post) => post.id !== postId));
          } else {
            const errorData = await response.json();
            alert(`Failed to delete post: ${errorData.message}`);
          }
        } catch (error) {
          console.error("Error deleting post:", error);
        }
    };

    useEffect(() => {
        fetchPosts();
    }, []);

    const handleViewProfileClick = () => {
        setShowViewProfilePage(true); // Show the ViewProfile page
        setShowNetworksPage(false); // Hide Networks page
        setSelectedUser(null); // Clear selected user
        setDropdownVisible(false); // Close the dropdown menu
    };

    // Handle clicking Home button
    const handleHomeClick = () => {
        setShowViewProfilePage(false);
        setShowNetworksPage(false);
        setSelectedUser(null);
        fetchPosts();
        fetchComments();
    };

    const handleLike = async (startupId) => {
        const isLiked = likedPosts.has(startupId);
        const currentCount = likeCounts[startupId] || 0;
      
        // Prevent unliking if like count is already 0 or less
        if (isLiked && currentCount < 0) return;
      
        const endpoint = isLiked ? 'unlike' : 'like';
        const method = isLiked ? 'DELETE' : 'POST';
      
        try {
          const response = await fetch(`http://localhost:8080/v1/startup/${endpoint}`, {
            method: method,
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              email: session?.user?.email, // user email from session
              startup_id: startupId
            }),
          });
      
          const data = await response.json();
      
          if (response.ok && data.status === "success") {
            // Update liked posts set
            if (session?.user?.email) {
                setLikedPosts((prev) => {
                const updated = new Set(prev);
                if (isLiked) {
                    updated.delete(startupId);
                } else {
                    updated.add(startupId);
                }
                localStorage.setItem("likedPosts", JSON.stringify(Array.from(updated)));
                return updated;
                });
        
                // Update like count
                setLikeCounts((prevCounts) => ({
                ...prevCounts,
                [startupId]: Math.max(0, (prevCounts[startupId] || 0) + (isLiked ? -1 : 1)),
                }));
            }
          } else {
            alert(data.message || `Failed to ${endpoint} post.`);
          }
        } catch (error) {
          console.error(`Error trying to ${endpoint} post:`, error);
        }
      };
            
    const fetchLikeCounts = async () => {
        const counts = {};
        for (const post of posts) {
          try {
            const response = await fetch(`http://localhost:8080/v1/startup/likes/${post.id}`);
            const data = await response.json();
            if (response.ok && data.status === "success") {
              counts[post.id] = data.like_count;
            } else {
              counts[post.id] = 0;
            }
          } catch {
            counts[post.id] = 0;
          }
        }
        setLikeCounts(counts);
    };

    const fetchCommentCounts = async () => {
        const counts = {};
        for (const post of posts) {
          try {
            const response = await fetch(`http://localhost:8080/v1/startup/comments/${post.id}`);
            const data = await response.json();
            if (response.ok && data.status === "success") {
              counts[post.id] = data.comments.length;
            } else {
              counts[post.id] = 0;
            }
          } catch {
            counts[post.id] = 0;
          }
        }
        setCommentCounts(counts);
      };

    useEffect(() => {
        const stored = localStorage.getItem("likedPosts");
        if (stored) {
          setLikedPosts(new Set(JSON.parse(stored)));
        }
      }, []);

    useEffect(() => {
        if (posts.length > 0) {
          fetchLikeCounts();
          fetchCommentCounts();
        }
      }, [posts]);
    
      const handleCommentClick = (postId) => {
        //setSelectedPostId(postId === selectedPostId ? null : postId);
        setOpenCommentPostId(openCommentPostId === postId ? null : postId);
      };
      
      const handleCommentSubmit = async (postId) => {
        if (!commentText.trim()) return;
      
        try {
          const response = await fetch('http://localhost:8080/v1/startup/comment', {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              email: session?.user?.email,
              startup_id: postId,
              content: commentText
            }),
          });
      
          const data = await response.json();
          if (response.ok && data.status === "success") {
            setCommentText('');
            fetchComments(postId); // Refresh comments
          } else {
            alert(data.message || "Failed to post comment");
          }
        } catch (error) {
          console.error("Error posting comment:", error);
        }
      };
      
    const fetchComments = async (postId) => {
        try {
          const response = await fetch(`http://localhost:8080/v1/startup/comments/${postId}`);
          const data = await response.json();
          if (response.ok && data.status === "success") {
            setComments(prev => ({
              ...prev,
              [postId]: data.comments
            }));
          }
        } catch (error) {
          console.error("Error fetching comments:", error);
        }
    };

    const handleShowComments = (postId) => {
        setVisibleCommentsPostId(visibleCommentsPostId === postId ? null : postId);
        if (!comments[postId]) fetchComments(postId); // Fetch if not cached
    };

    const handleReportSubmit = async () => {
        if (!reportReason.trim()) {
          alert("Please provide a reason for reporting.");
          return;
        }
      
        try {
          const response = await fetch('http://localhost:8080/v1/startup/report', {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              email: session?.user?.email,
              startup_id: reportModalPostId,
              reason: reportReason
            }),
          });
      
          const data = await response.json();
          if (response.ok && data.status === "success") {
            alert(data.message);
            setReportModalPostId(null);
            setReportReason('');
            fetchPosts(); // Refresh posts if a startup was deleted
          } else {
            alert(data.message || "Failed to submit report");
          }
        } catch (error) {
          console.error("Error submitting report:", error);
          alert("Error submitting report");
        }
      };
      

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
                    <button className={`${styles.navItem} ${styles.active}`} onClick={handleHomeClick}>
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="none">
                        <path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z" fill="#00DC82"/>
                        </svg>
                        Home</button>
                    <button className={`${styles.navItem} ${styles.active}`} onClick={handleMyNetworkClick} >
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="none">
                        <path d="M12 12c2.7 0 5-2.3 5-5-2.3-5-5-5-5 0s2.3 5 5 5zm0 2c-3.9 0-7.2 2.1-9 5.2V22h18v-2.8c-1.8-3.1-5.1-5.2-9-5.2z" fill="#00DC82"/>
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
                            <button className={styles.dropdownItem} onClick={handleViewProfileClick}>View Profile</button>
                            <button className={styles.dropdownItem} onClick={() => signOut({ callbackUrl: "/" })}>Sign Out</button>
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
            
            {/* Main Content Rendering Logic */}
            {showViewProfilePage ? (
                <ViewProfile />
            ) : showNetworksPage ? (
                <Networks />
            ) : (
                <>
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
                                                        <button className={styles.optionButton} onClick={() => handleEditClick(post)}>
                                                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                            <path d="M12 20h9"/>
                                                            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4Z"/>
                                                            </svg>
                                                        Edit
                                                        </button>
                                                        <button className={styles.optionButton} onClick={() => handleDeletePost(post.id)}>
                                                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                                <path d="M3 6h18"/>
                                                                <path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                                                                <path d="M10 11v6"/>
                                                                <path d="M14 11v6"/>
                                                                <path d="M5 6h14l-1 14H6Z"/>
                                                            </svg>
                                                            Delete</button>
                                                        <button className={styles.optionButton} onClick={() => setReportModalPostId(post.id)}>
                                                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                            <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
                                                            <line x1="12" y1="9" x2="12" y2="13"/>
                                                            <line x1="12" y1="17" x2="12.01" y2="17"/>
                                                            </svg>
                                                            Report</button>
                                                    </div>
                                                )}
                                                {reportModalPostId && (
                                                    <div className={styles.reportOverlay}>
                                                        <div className={styles.reportContent}>
                                                        <h3>Report this post</h3><br></br>
                                                        <textarea
                                                            value={reportReason}
                                                            onChange={(e) => setReportReason(e.target.value)}
                                                            placeholder="Please specify the reason for reporting..."
                                                            rows={4}
                                                            className={styles.reportText}
                                                        />
                                                        <div className={styles.reportActions}>
                                                            <button 
                                                            onClick={() => handleReportSubmit(reportModalPostId)}
                                                            disabled={!reportReason.trim()}
                                                            className={styles.reportButton}
                                                            >
                                                            Submit
                                                            </button>
                                                            <button onClick={() => {
                                                            setReportModalPostId(null);
                                                            setReportReason('');
                                                            }}
                                                            className={styles.reportButton}>
                                                            Cancel
                                                            </button>
                                                        </div>
                                                        </div>
                                                    </div>
                                                    )}

                                                {editingPostId === post.id ? (
                                                    <>
                                                        <div className={styles.editPost}>
                                                            <label>
                                                            Startup Name:
                                                            <input
                                                                type="text"
                                                                className={styles.formStyle}
                                                                value={editedPostData.startup_name}
                                                                onChange={(e) => handleEditChange(e, "startup_name")}
                                                            />
                                                            </label><br></br><br></br>
                                                            <label>
                                                            Industry:<br></br>
                                                            <input
                                                                type="text"
                                                                className={styles.formStyle}
                                                                value={editedPostData.industry}
                                                                onChange={(e) => handleEditChange(e, "industry")}
                                                            />
                                                            </label><br></br><br></br>
                                                            <label>
                                                            Description:
                                                            <textarea
                                                                className={styles.formStyle}
                                                                value={editedPostData.description}
                                                                onChange={(e) => handleEditChange(e, "description")}
                                                            />
                                                            </label><br></br><br></br>
                                                            <label>
                                                            Budget:<br></br>
                                                            <input
                                                                type="text"
                                                                className={styles.formStyle}
                                                                value={editedPostData.budget}
                                                                onChange={(e) => handleEditChange(e, "budget")}
                                                            />
                                                            </label><br></br><br></br>
                                                            <label>
                                                            Timeframe:
                                                            <input
                                                                type="text"
                                                                className={styles.formStyle}
                                                                value={editedPostData.timeframe}
                                                                onChange={(e) => handleEditChange(e, "timeframe")}
                                                            />
                                                            </label><br></br><br></br>
                                                        {/* Save and Cancel Buttons */}
                                                            <div className={styles.editButtons}>
                                                            <button className={styles.likeButton} onClick={() => handleSaveEdit(post.id)}>Save</button>
                                                            <button className={styles.commentButton} onClick={() => setEditingPostId(null)}>Cancel</button>
                                                            </div>
                                                        </div>
                                                    </>
                                                ) : (
                                                    <>
                                                    <h3 className={styles.postStyle}>{post.entrepreneur.name}</h3>
                                                    <p className={styles.timeStyle}>
                                                        {moment(post.created_at).fromNow()}
                                                    </p>
                                                    {post.updated_at && post.updated_at !== post.created_at && (
                                                        <p className={styles.timeStyle}>
                                                        {moment(post.updated_at).fromNow()}
                                                        </p>
                                                    )} <br></br>
                                                    <h4 className={styles.postStyle}>Startup Name: {post.startup_name}</h4>
                                                    <p className={styles.postStyle}>Industry: {post.industry}</p>
                                                    <p className={styles.postStyle}>Description: {post.description}</p>
                                                    <p className={styles.postStyle}>Budget: {post.budget}</p>
                                                    <p className={styles.postStyle}>Timeframe: {post.timeframe}</p>
                                                    <hr className={styles.separator}></hr>
                                                    <div className={styles.buttonContainer}>
                                                        <div className={styles.likeWrapper}>
                                                            <div className={styles.likeCount}>
                                                            {likeCounts[post.id] !== undefined ? `${likeCounts[post.id]} Like(s)` : "Loading..."}
                                                            </div>
                                                            <button
                                                            className={`${styles.likeButton} ${likedPosts.has(post.id) ? styles.liked : styles.likeButton}`}
                                                            onClick={() => handleLike(post.id)}
                                                            >
                                                            Like
                                                            </button>
                                                        </div>
                                                        <div className={styles.likeWrapper}>
                                                            <div className={styles.likeCount}>
                                                            <button 
                                                                onClick={() => handleShowComments(post.id)}
                                                                className={styles.commentCountButton}
                                                                >
                                                                {commentCounts[post.id] || 0} Comment(s)
                                                                </button>
                                                            </div>
                                                            <button className={styles.commentButton} onClick={() => handleCommentClick(post.id)}>Comment</button>
                                                        </div>
                                                    </div>
                                                    {openCommentPostId === post.id && (
                                                        <div className={styles.commentSection}>
                                                            <textarea
                                                                value={commentText}
                                                                onChange={(e) => setCommentText(e.target.value)}
                                                                onKeyDown={(e) => {
                                                                    if (e.key === 'Enter' && !e.shiftKey) {
                                                                    e.preventDefault();
                                                                    handleCommentSubmit(post.id);
                                                                    }
                                                                }}
                                                                placeholder="Add a comment..."
                                                                className={styles.commentInput}
                                                            />
                                                        </div>
                                                    )}
                                                    {visibleCommentsPostId === post.id && comments[post.id] && (
                                                        <div>
                                                            {comments[post.id]?.map((comment) => (
                                                            <div key={comment.id} className={styles.comment}>
                                                                <strong>{comment.user_name}</strong><br></br>
                                                                <span className={styles.timeDisplay}>{moment(comment.created_at).fromNow()}</span>
                                                                <p>{comment.content}</p>
                                                            </div>
                                                        ))}
                                                        </div>
                                                    )}    
                                                    </>
                                                )}
                                            </div>
                                    ))}
                                </div>
                            </div>
                        </div>
                        </>
                    )}
                </div>
                </>
            )}
            
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