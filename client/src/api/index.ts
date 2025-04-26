import ky from "ky";
import { isValid, parseISO, isDate } from "date-fns";
import { getTimezoneOffset } from "date-fns-tz";

export const api = ky.create({
  prefixUrl: import.meta.env.VITE_API_URL,
  // credentials: "include",
  headers: {
    "Content-Type": "application/json",
  },
  hooks: {
    beforeRequest: [
      (request, opt: any) => {
        const token = localStorage.getItem(`@${import.meta.env.VITE_APP_NAME}:token`);
        request.headers.set("Authorization", `Bearer ${token}`);

        if (!opt.searchParams) return request;

        const keys = Object.keys(opt.searchParams);
        let requestUrl = request.url;

        for (const key of keys) {
          const val = opt.searchParams[key];
          if (val && val !== "") {
            let validDate = false;
            try {
              const dt = parseISO(val);
              validDate = isDate(dt) && isValid(dt);

              if (validDate) {
                const tz = Intl.DateTimeFormat().resolvedOptions().timeZone;
                const tzOffset = getTimezoneOffset(tz, dt);
                dt.setMilliseconds(tzOffset);

                requestUrl = requestUrl.replace(encodeURIComponent(val), encodeURIComponent(dt.toISOString()));
                opt.searchParams[key] = dt.toISOString();
              }
            } catch {}
          }
        }

        const newReq = new Request(requestUrl, request);

        //console.log({ opt, request, newReq });

        return newReq;
      },
    ],
    // afterResponse: [
    //   async (request, _, response) => {
    //     if (!response.ok) {
    //       if (response.status === 401 && !request.url.includes("/auth/login") && !request.url.includes("/auth/admin")) {
    //         localStorage.removeItem(`@${import.meta.env.VITE_APP_NAME}:token`);
    //         window.location.href = "#/login";
    //         return;
    //       }

    //       throw new Error(await response.json());
    //     }

    //     return response;
    //   },
    // ],
  },
});
