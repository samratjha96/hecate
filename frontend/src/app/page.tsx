"use client";

import SubredditDashboard from "@/components/SubredditDashboard";
import SearchBox from "@/components/SearchBox";

export default function Home() {
  return (
    <div className="min-h-screen p-4 sm:p-8 font-[family-name:var(--font-geist-sans)]">
      <main>
        <SearchBox />
        <SubredditDashboard />
      </main>
    </div>
  );
}
