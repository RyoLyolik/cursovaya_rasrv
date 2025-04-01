'use server';

import axios from "axios";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const sessionCookie = request.cookies.get('brick-session-id')
  try {
    const resp = await axios.get(
        'http://localhost:8080/account/me',
        {
            headers: {
                Cookie: `${sessionCookie?.name}=${sessionCookie?.value}`
            },
            withCredentials: true
        }
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