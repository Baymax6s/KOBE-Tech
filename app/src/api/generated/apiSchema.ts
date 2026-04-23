/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface ArticleArticleJSON {
  /** @example "Article body" */
  content?: string;
  /** @example "2026-04-24T00:00:00Z" */
  created_at?: string;
  /** @example 1 */
  id?: number;
  /** @example "First article" */
  title?: string;
  /** @example "2026-04-24T00:00:00Z" */
  updated_at?: string;
  /** @example 1 */
  user_id?: number;
}

export interface ServerArticleErrorResponse {
  /** @example "internal server error" */
  message?: string;
}

export interface ServerCreateArticleRequest {
  /** @example "Article body" */
  body?: string;
  /** @example "First article" */
  title?: string;
}

export interface ServerListArticlesResponse {
  articles?: ArticleArticleJSON[];
}

export interface ServerLoginRequest {
  /**
   * @format email
   * @example "user@example.com"
   */
  email?: string;
  /**
   * @format password
   * @example "change-me"
   */
  password?: string;
}

export interface ServerNotImplementedResponse {
  /** @example "feature is not implemented yet" */
  message?: string;
  /** @example "internal/{domain}/{handler,service,repository}.go" */
  next_step?: string;
}

import type {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  HeadersDefaults,
  ResponseType,
} from "axios";
import axios from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams
  extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown>
  extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({
    securityWorker,
    secure,
    format,
    ...axiosConfig
  }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({
      ...axiosConfig,
      baseURL: axiosConfig.baseURL || "/",
    });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected mergeRequestParams(
    params1: AxiosRequestConfig,
    params2?: AxiosRequestConfig,
  ): AxiosRequestConfig {
    const method = params1.method || (params2 && params2.method);

    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...((method &&
          this.instance.defaults.headers[
            method.toLowerCase() as keyof HeadersDefaults
          ]) ||
          {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected stringifyFormItem(formItem: unknown) {
    if (typeof formItem === "object" && formItem !== null) {
      return JSON.stringify(formItem);
    } else {
      return `${formItem}`;
    }
  }

  protected createFormData(input: Record<string, unknown>): FormData {
    if (input instanceof FormData) {
      return input;
    }
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      const propertyContent: any[] =
        property instanceof Array ? property : [property];

      for (const formItem of propertyContent) {
        const isFileType = formItem instanceof Blob || formItem instanceof File;
        formData.append(
          key,
          isFileType ? formItem : this.stringifyFormItem(formItem),
        );
      }

      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = format || this.format || undefined;

    if (
      type === ContentType.FormData &&
      body &&
      body !== null &&
      typeof body === "object"
    ) {
      body = this.createFormData(body as Record<string, unknown>);
    }

    if (
      type === ContentType.Text &&
      body &&
      body !== null &&
      typeof body !== "string"
    ) {
      body = JSON.stringify(body);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type ? { "Content-Type": type } : {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title KOBE-Tech API
 * @version 0.1.0
 * @baseUrl /
 * @contact
 *
 * KOBE-Tech API documentation generated from Go annotations.
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  api = {
    /**
     * @description Get article list API.
     *
     * @tags articles
     * @name ArticlesList
     * @summary List articles
     * @request GET:/api/articles
     */
    articlesList: (params: RequestParams = {}) =>
      this.request<ServerListArticlesResponse, ServerArticleErrorResponse>({
        path: `/api/articles`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description Create article API.
     *
     * @tags articles
     * @name ArticlesCreate
     * @summary Create article
     * @request POST:/api/articles
     */
    articlesCreate: (
      request: ServerCreateArticleRequest,
      params: RequestParams = {},
    ) =>
      this.request<any, ServerNotImplementedResponse>({
        path: `/api/articles`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description Login API.
     *
     * @tags auth
     * @name AuthLoginCreate
     * @summary Login
     * @request POST:/api/auth/login
     */
    authLoginCreate: (
      request: ServerLoginRequest,
      params: RequestParams = {},
    ) =>
      this.request<any, ServerNotImplementedResponse>({
        path: `/api/auth/login`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        ...params,
      }),
  };
}
