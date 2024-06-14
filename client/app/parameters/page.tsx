/**
 * eslint-disable @next/next/no-img-element
 *
 * @format
 */

/**
 * eslint-disable @next/next/no-img-element
 *
 * @format
 */

/** @format */
"use client";

import { DataTable } from "@/components/DataTable";
import React from "react";
import PageTitle from "@/components/PageTitle";
import { cn } from "@/lib/utils";

type Props = {};

export default function OrdersPage({}: Props) {
    return (
        <div className="flex flex-col gap-5  w-full">
            <PageTitle
                title="Parameters"
                className="pb-10 text-5xl font-normal"
            />
        </div>
    );
}
