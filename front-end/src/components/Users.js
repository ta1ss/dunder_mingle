import { useEffect, useState } from "react";
import { Link, useOutletContext } from "react-router-dom";

const Users = () => {

    const { userId } = useOutletContext();
    const [allUsers, setAllUsers] = useState([]);
    const [searchQuery, setSearchQuery] = useState('');
    const currentUser = allUsers.find((user) => user.id == Number(userId));
    const [loading, setLoading] = useState(true);
    const {setAlertClassName} = useOutletContext();

    const profileImagesEndpoint = "http://localhost:8080/media/profile_images/"

    const fetchAllUsers = () => {
        fetch("http://localhost:8080/users", {
            method: "GET",
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                setAllUsers(data);
                setLoading(false);
                // console.log(data)

            })
            .catch((error) => {
                console.log(error);
            });
    }

    const sortedUsers = allUsers.sort((a, b) => {
        if (a.firstName < b.firstName) {
            return -1;
        }
        if (a.firstName > b.firstName) {
            return 1;
        }
        return 0;
    })

    const filteredUsers = sortedUsers.filter((user) => {
        if (currentUser !== undefined) {
            const fullName = `${user.firstName} ${user.lastName}`.toLowerCase();
            return fullName.includes(searchQuery.toLowerCase()) && user.id !== currentUser.id;
        }
    });

    useEffect(() => {
        fetchAllUsers();
        setAlertClassName('d-none');
    }, []);

    return (
        <>
            {loading ? ("") : (
                <div className="container users-container">
                    <div className="row">
                        <div className="users-search">
                            <input
                                type="text"
                                placeholder="Search users..."
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                className="form-control mb-3 users-search-input"
                            />
                        </div>
                        <div className="col-md-6">
                            {filteredUsers
                                .slice(0, Math.ceil((filteredUsers.length) / 2))
                                .map((user) => (
                                    <div key={user.id}>
                                        <Link to={`/profile/${user.id}`} className="users-users-link">
                                            <div className='d-flex users-user-container' >
                                                <img src={`${profileImagesEndpoint}${user.image}`} alt="profile image" className="users-user-image img-fluid" />
                                                <p className="users-user-name">{user.firstName} {user.lastName} </p>
                                            </div>
                                        </Link>
                                    </div>
                                ))}
                        </div>
                        <div className="col-md-6">
                            {filteredUsers
                                .slice(Math.ceil(filteredUsers.length / 2))
                                .map((user) => (
                                    <div key={user.id}>
                                        <Link to={`/profile/${user.id}`} className="users-users-link">
                                            <div className='d-flex users-user-container' >
                                                <img src={`${profileImagesEndpoint}${user.image}`} alt="profile image" className="users-user-image img-fluid" />
                                                <p className="users-user-name">{user.firstName} {user.lastName}</p>
                                            </div>
                                        </Link>
                                    </div>
                                ))}
                        </div>
                    </div>
                </div>
            )}
        </>
    )
}

export default Users;