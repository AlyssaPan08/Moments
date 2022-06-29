import React, {useEffect, useState} from 'react';
import { Tabs, message, Row, Col} from "antd";
import axios from "axios";
import SearchBar from './SearchBar';
import CreatePostButton from "./CreatePostButton";
import { SEARCH_KEY, BASE_URL, TOKEN_KEY} from "../constants";
import PhotoGallery from "./PhotoGallery";

const { TabPane } = Tabs;

function Home(props) {
    //state: posts, activeTab, searchOption
    const [posts, setPost] = useState([]);
    const [activeTab, setActiveTab] = useState("image");
    const [searchOption, setSearchOption] = useState( {
        type: SEARCH_KEY.all,
        keyword: ""
    })

    const handleSearch = (option) => {
        const { type, keyword } = option;
        setSearchOption({ type: type, keyword: keyword });
    };

    //component life cycle
    //case 1: do search at first time -> did mount -> search {type: all, keyword:""}
    //case 2: di search after the first time -> did update -> search {type: keyword/user, keyword:value}
    useEffect( () => {
        fetchPost(searchOption)
    }, [searchOption]);

    const fetchPost = option => {
        //step 1: get search type & search context
        //step 2: fetch post from server
        //step 3: analyze response from server
        //  case 1: success -> display posts (image/video)
        //  case 2: fail -> inform user

        const { type, keyword } = option; //step 1

        //get url based on option
        let url = "";
        if (type === SEARCH_KEY.all) {
            url = `${BASE_URL}/search`;
        } else if (type === SEARCH_KEY.user) {
            url = `${BASE_URL}/search?user=${keyword}`;
        } else {
            url = `${BASE_URL}/search?keywords=${keyword}`;
        }

        //configure
        const opt = {
            method: "GET",
            url: url,
            headers: {
                Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`
            }
        }

        axios(opt)
            .then((res) => {
                if (res.status === 200) {
                    setPost(res.data);
                }
            })
            .catch((err) => {
                message.error("Fetch posts failed!");
                console.log("fetch posts failed: ", err.message);
            });

    }

    const renderPosts = (type) => {
        //case 1: no any post
        //case 2: type = "image" => display all images (use Gallery)
        //cas3 2: type == "video" => display all videos
        if (!posts || posts.length === 0) {
            return <div>No data!</div>;
        }

        if (type === "image") {
            //render image
            const imageArr = posts
                .filter((item) => item.type === "image")
                .map((image) => {
                    return {
                        postId: image.id,
                        src: image.url,
                        user: image.user,
                        caption: image.message,
                        thumbnail: image.url,
                        thumbnailWidth: 300,
                        thumbnailHeight: 200
                    };
                });
            return <PhotoGallery images={imageArr} />;
        } else if (type === "video") {
            return (
                <Row gutter={32}>
                    {posts
                        .filter((post) => post.type === "video")
                        .map((post) => (
                            <Col span={8} key={post.url}>
                                <video src={post.url} controls={true} className="video-block" />
                                <p>
                                    {post.user}: {post.message}
                                </p>
                            </Col>
                        ))}
                </Row>
            );
        }
    };

    const showPost = (type) => {
        console.log("type -> ", type);
        setActiveTab(type);

        setTimeout(() => {
            setSearchOption({ type: SEARCH_KEY.all, keyword: "" });
        }, 3000);
    };

    const operations = <CreatePostButton onShowPost={showPost} />;

    return (
        <div className="home">
            <SearchBar handleSearch={handleSearch}/>
            <div className="display">
                <Tabs
                    onChange={(key) => setActiveTab(key)}
                    defaultActiveKey="image"
                    activeKey={activeTab}
                    tabBarExtraContent={operations}
                >
                    <TabPane tab="Images" key="image">
                        {renderPosts("image")}
                    </TabPane>
                    <TabPane tab="Videos" key="video">
                        {renderPosts("video")}
                    </TabPane>
                </Tabs>
            </div>
        </div>
    );
}

export default Home;