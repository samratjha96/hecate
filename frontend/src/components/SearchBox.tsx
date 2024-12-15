import React, { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { API_BASE_URL } from "@/lib/utils";
import { toast } from "sonner";

interface SearchResult {
  title: string;
  content: string;
  discussionUrl: string;
  commentCount: number;
  upvotes: number;
  subredditName: string;
}

interface SearchResponse {
  posts: SearchResult[];
}

const SearchBox = () => {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(false);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;

    setLoading(true);
    try {
      const response = await fetch(
        `${API_BASE_URL}/subreddits/search?q=${encodeURIComponent(query)}`
      );
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data: SearchResponse = await response.json();
      setResults(data.posts);
    } catch (error) {
      console.error("Search failed:", error);
      toast.error("Failed to search posts");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="mb-8">
      <CardContent className="pt-6">
        <form onSubmit={handleSearch} className="flex gap-4 mb-6">
          <Input
            type="text"
            placeholder="Search across all posts..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="flex-grow"
          />
          <Button type="submit" disabled={loading}>
            {loading ? "Searching..." : "Search"}
          </Button>
        </form>

        {results.length > 0 && (
          <div className="space-y-4">
            <h3 className="text-lg font-semibold">Search Results</h3>
            <div className="divide-y">
              {results.map((result, index) => (
                <div key={index} className="py-4">
                  <Button
                    variant="link"
                    className="p-0 h-auto text-left font-semibold hover:no-underline"
                    onClick={() => window.open(result.discussionUrl, "_blank")}
                  >
                    {result.title}
                  </Button>
                  <div className="text-sm text-muted-foreground mt-1">
                    r/{result.subredditName} • {result.upvotes} upvotes •{" "}
                    {result.commentCount} comments
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default SearchBox;
