/** @format */

import PageTitle from "@/components/PageTitle";
import Image from "next/image";
import { DollarSign, Users, CreditCard, Activity } from "lucide-react";
import Card, { CardProps } from "@/components/Card";

export default function Home() {
    return (
        <div className="flex flex-col gap-5  w-full">
            <PageTitle title="CMS" className="pb-10 text-5xl font-normal" />
        </div>
    );
}
