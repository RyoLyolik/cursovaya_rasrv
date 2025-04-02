'use server';

import axios from "axios";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  try {
    const resp = await axios.get(
        `http://localhost:9999/api/report/file/list`,
    );
    const response = NextResponse.json(resp.data);
    return response;
  } catch (err) {
    if (axios.isAxiosError(err)) {
      return NextResponse.json(err.response?.data || {status: 'Error'});
    } else {
      return NextResponse.json({status: 'Error'});
    }
  }
}