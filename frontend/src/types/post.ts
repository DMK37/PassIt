export type Post = {
    id: string;
    user_id: string;
    user: PostUser;
    text: string;
    images: string[];
    timestamp: number;
};

export type PostUser = {
    username: string;
    first_name: string;
    last_name: string;
    avatar: string;
}
