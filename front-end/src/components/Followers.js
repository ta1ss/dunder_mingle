import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

const Followers = (props) => {

  const [data, setData] = useState({ followers: [], following: [] })
  const profileImagesEndpoint = "http://localhost:8080/media/profile_images/"

  useEffect(() => {
    const options = {
      method: 'GET',
      credentials: 'include',
    }

    const followersEndpoint = 'http://localhost:8080/profile/followers' + (props.userId ? `?userId=${props.userId}` : '')
    fetch(followersEndpoint, options)
      .then(response => response.json())
      .then((data) => {
        const followers = data.followers ? data.followers.map(follower => ({ followerId: follower.followerId, followerName: follower.followerName, followerImg: follower.userImg })) : [];
        const following = data.following ? data.following.map(following => ({ followingId: following.followingId, followingName: following.followingName, followingImg: following.userImg })) : [];
        setData({ followers, following })
      })
      .catch(error => console.log(error))


  }, [])


  return (
    <div className="container">
      <div className="row">
        <div className="col-md-6 followers-border-right followers-container">
          <h3 className="followers-h3">Followers ({data.followers.length})</h3>
          {data.followers && data.followers.length > 0 ? (
            <div className="row">
              {data.followers.map((follower, index) => (
                <div key={index} className="follower-box">
                  <Link to={`/profile/${follower.followerId}`} className="users-users-link">
                    <div className='profile-followers' >
                      <img src={`${profileImagesEndpoint}${follower.followerImg}`} alt="profile image" className="profile-follower-image img-fluid" />
                      <p className="profile-follower-name">{follower.followerName} </p>
                    </div>
                  </Link>
                </div>
              ))}
            </div>
          ) : (
            <p className="followers-follower">This user has no followers.</p>
          )}
        </div>
        <div className="col-md-6 followers-container">
          <div>
            <h3 className="followers-h3">Following ({data.following.length})</h3>
            {data.following && data.following.length > 0 ? (
              <div className="row">
                  {data.following.map((following, index) => (
                    <div key={index} className="following-box">
                      <Link to={`/profile/${following.followingId}`} className="users-users-link">
                        <div className='profile-followers' >
                          <img src={`${profileImagesEndpoint}${following.followingImg}`} alt="profile image" className="profile-follower-image img-fluid" />
                          <p className="profile-follower-name">{following.followingName} </p>
                        </div>
                      </Link>
                    </div>
                  ))}
                </div>
            ) : (
              <p className="followers-follower">This user is not following anyone.</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );

};

export default Followers;