/** @format */

import PageTitle from "@/components/PageTitle";
import Image from "next/image";
import Card, { CardProps } from "@/components/Card";

export default function Home() {
  return (
    <div className="flex flex-col gap-5  w-full">
      <PageTitle title="Library" />
    </div>
  );
}
