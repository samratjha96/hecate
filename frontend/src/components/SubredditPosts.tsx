import React, { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { API_BASE_URL } from "@/lib/utils";

interface Post {
  title: string;
  content: string;
  discussionUrl: string;
  commentCount: number;
  upvotes: number;
}

const SubredditPosts = () => {
  const [posts, setPosts] = useState<Post[]>([]);
  const [subredditName, setSubredditName] = useState("");
  const [loading, setLoading] = useState(false);

  const fetchPosts = async (subreddit: string) => {
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE_URL}/subreddits/${subreddit}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setPosts(data);
    } catch (error) {
      console.error("Error fetching posts:", error);
      toast.error(`Failed to fetch posts from r/${subreddit}`);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (subredditName) {
      fetchPosts(subredditName);
    }
  };

  return (
    <div className="container mx-auto p-4">
      <Card className="mb-8">
        <CardHeader>
          <CardTitle className="flex justify-between items-center">
            Subreddit Posts
            <form onSubmit={handleSubmit} className="flex space-x-2">
              <Input
                type="text"
                placeholder="Enter subreddit name"
                value={subredditName}
                onChange={(e) => setSubredditName(e.target.value)}
                className="max-w-xs"
              />
              <Button type="submit" disabled={loading}>
                {loading ? "Loading..." : "Fetch Posts"}
              </Button>
            </form>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {posts && posts.length > 0 ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Title</TableHead>
                  <TableHead>Upvotes</TableHead>
                  <TableHead>Comments</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {posts.map((post) => (
                  <TableRow key={post.title}>
                    <TableCell className="max-w-4xl truncate overflow-hidden">
                        <Button variant="link" onClick={() => window.open(post.discussionUrl, "_blank")}>
                            {post.title}
                        </Button>
                    </TableCell>
                    <TableCell>{post.upvotes}</TableCell>
                    <TableCell>{post.commentCount}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : (
            <p>Nothing to display. Enter a subreddit name to view posts</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default SubredditPosts;