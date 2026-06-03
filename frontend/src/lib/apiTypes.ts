/**
 * TypeScript types matching contracts/calculate-openapi.yaml.
 */

export interface CalculationRequest {
  expression: string;
}

export interface CalculationResult {
  result: string;
  expression: string;
}

export interface ApiErrorDetail {
  code: string;
  message: string;
  fields?: Record<string, string>;
}

export interface ApiError {
  error: ApiErrorDetail;
}

/** Thrown by calculate() on non-2xx responses. */
export class ApiResponseError extends Error {
  readonly code: string;
  readonly fields?: Record<string, string>;

  constructor(code: string, message: string, fields?: Record<string, string>) {
    super(message);
    this.name = "ApiResponseError";
    this.code = code;
    this.fields = fields;
  }
}
