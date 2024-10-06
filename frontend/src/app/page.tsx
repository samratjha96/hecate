"use client";

import SubredditDashboard from "@/components/SubredditDashboard";
import SubredditPosts from "@/components/SubredditPosts";

export default function Home() {
  return (
    <div className="flex items-center justify-items-center p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
          <div className="flex flex-col md:flex-row">
            <SubredditDashboard />
            <SubredditPosts />
          </div>
      </main>
    </div>
  );
}
