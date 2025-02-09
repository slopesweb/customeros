// @ts-nocheck remove this when typscript-react-query plugin is fixed
import * as Types from '../../../src/types/__generated__/graphql.types';

import { GraphQLClient } from 'graphql-request';
import { RequestInit } from 'graphql-request/dist/types.dom';
import {
  useQuery,
  useInfiniteQuery,
  UseQueryOptions,
  UseInfiniteQueryOptions,
  InfiniteData,
} from '@tanstack/react-query';

function fetcher<TData, TVariables extends { [key: string]: any }>(
  client: GraphQLClient,
  query: string,
  variables?: TVariables,
  requestHeaders?: RequestInit['headers'],
) {
  return async (): Promise<TData> =>
    client.request({
      document: query,
      variables,
      requestHeaders,
    });
}
export type GetRenewalsQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
  where?: Types.InputMaybe<Types.Filter>;
  sort?: Types.InputMaybe<Types.SortBy>;
}>;

export type GetRenewalsQuery = {
  __typename?: 'Query';
  dashboardView_Renewals?: {
    __typename?: 'RenewalsPage';
    totalPages: number;
    totalElements: any;
    totalAvailable: any;
    content: Array<{
      __typename?: 'RenewalRecord';
      organization: {
        __typename?: 'Organization';
        name: string;
        logo?: string | null;
        metadata: { __typename?: 'Metadata'; id: string };
        owner?: {
          __typename?: 'User';
          name?: string | null;
          firstName: string;
          lastName: string;
        } | null;
        accountDetails?: {
          __typename?: 'OrgAccountDetails';
          renewalSummary?: {
            __typename?: 'RenewalSummary';
            renewalLikelihood?: Types.OpportunityRenewalLikelihood | null;
            nextRenewalDate?: any | null;
            arrForecast?: number | null;
            maxArrForecast?: number | null;
          } | null;
        } | null;
        lastTouchpoint?: {
          __typename?: 'LastTouchpoint';
          lastTouchPointAt?: any | null;
          lastTouchPointType?: Types.LastTouchpointType | null;
          lastTouchPointTimelineEvent?:
            | {
                __typename: 'Action';
                id: string;
                actionType: Types.ActionType;
                createdAt: any;
                source: Types.DataSource;
                createdBy?: {
                  __typename?: 'User';
                  id: string;
                  firstName: string;
                  lastName: string;
                } | null;
              }
            | { __typename: 'Analysis'; id: string }
            | {
                __typename: 'InteractionEvent';
                id: string;
                channel?: string | null;
                eventType?: string | null;
                externalLinks: Array<{
                  __typename?: 'ExternalSystem';
                  type: Types.ExternalSystemType;
                }>;
                sentBy: Array<
                  | {
                      __typename: 'ContactParticipant';
                      contactParticipant: {
                        __typename?: 'Contact';
                        id: string;
                        name?: string | null;
                        firstName?: string | null;
                        lastName?: string | null;
                      };
                    }
                  | {
                      __typename: 'EmailParticipant';
                      type?: string | null;
                      emailParticipant: {
                        __typename?: 'Email';
                        id: string;
                        email?: string | null;
                        rawEmail?: string | null;
                      };
                    }
                  | {
                      __typename: 'JobRoleParticipant';
                      jobRoleParticipant: {
                        __typename?: 'JobRole';
                        contact?: {
                          __typename?: 'Contact';
                          id: string;
                          name?: string | null;
                          firstName?: string | null;
                          lastName?: string | null;
                        } | null;
                      };
                    }
                  | { __typename: 'OrganizationParticipant' }
                  | { __typename: 'PhoneNumberParticipant' }
                  | {
                      __typename: 'UserParticipant';
                      userParticipant: {
                        __typename?: 'User';
                        id: string;
                        firstName: string;
                        lastName: string;
                      };
                    }
                >;
              }
            | { __typename: 'InteractionSession' }
            | {
                __typename: 'Issue';
                id: string;
                createdAt: any;
                updatedAt: any;
              }
            | {
                __typename: 'LogEntry';
                id: string;
                createdBy?: {
                  __typename?: 'User';
                  lastName: string;
                  firstName: string;
                } | null;
              }
            | {
                __typename: 'Meeting';
                id: string;
                name?: string | null;
                attendedBy: Array<
                  | { __typename: 'ContactParticipant' }
                  | { __typename: 'EmailParticipant' }
                  | { __typename: 'OrganizationParticipant' }
                  | { __typename: 'UserParticipant' }
                >;
              }
            | {
                __typename: 'Note';
                id: string;
                createdBy?: {
                  __typename?: 'User';
                  firstName: string;
                  lastName: string;
                } | null;
              }
            | { __typename: 'Order' }
            | { __typename: 'PageView'; id: string }
            | null;
        } | null;
      };
      contract: {
        __typename?: 'Contract';
        contractName: string;
        metadata: { __typename?: 'Metadata'; id: string };
      };
      opportunity?: {
        __typename?: 'Opportunity';
        owner?: {
          __typename?: 'User';
          id: string;
          firstName: string;
          lastName: string;
          name?: string | null;
        } | null;
      } | null;
    }>;
  } | null;
};

export const GetRenewalsDocument = `
    query getRenewals($pagination: Pagination!, $where: Filter, $sort: SortBy) {
  dashboardView_Renewals(pagination: $pagination, where: $where, sort: $sort) {
    content {
      organization {
        metadata {
          id
        }
        name
        logo
        owner {
          name
          firstName
          lastName
        }
        accountDetails {
          renewalSummary {
            renewalLikelihood
            nextRenewalDate
            arrForecast
            maxArrForecast
          }
        }
        lastTouchpoint {
          lastTouchPointAt
          lastTouchPointType
          lastTouchPointTimelineEvent {
            __typename
            ... on PageView {
              id
            }
            ... on Issue {
              id
              createdAt
              updatedAt
            }
            ... on LogEntry {
              id
              createdBy {
                lastName
                firstName
              }
            }
            ... on Note {
              id
              createdBy {
                firstName
                lastName
              }
            }
            ... on InteractionEvent {
              id
              channel
              eventType
              externalLinks {
                type
              }
              sentBy {
                __typename
                ... on EmailParticipant {
                  type
                  emailParticipant {
                    id
                    email
                    rawEmail
                  }
                }
                ... on ContactParticipant {
                  contactParticipant {
                    id
                    name
                    firstName
                    lastName
                  }
                }
                ... on JobRoleParticipant {
                  jobRoleParticipant {
                    contact {
                      id
                      name
                      firstName
                      lastName
                    }
                  }
                }
                ... on UserParticipant {
                  userParticipant {
                    id
                    firstName
                    lastName
                  }
                }
              }
            }
            ... on Analysis {
              id
            }
            ... on Meeting {
              id
              name
              attendedBy {
                __typename
              }
            }
            ... on Action {
              id
              actionType
              createdAt
              source
              createdBy {
                id
                firstName
                lastName
              }
            }
          }
        }
      }
      contract {
        metadata {
          id
        }
        contractName
      }
      opportunity {
        owner {
          id
          firstName
          lastName
          name
        }
      }
    }
    totalPages
    totalElements
    totalAvailable
  }
}
    `;

export const useGetRenewalsQuery = <TData = GetRenewalsQuery, TError = unknown>(
  client: GraphQLClient,
  variables: GetRenewalsQueryVariables,
  options?: Omit<
    UseQueryOptions<GetRenewalsQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseQueryOptions<GetRenewalsQuery, TError, TData>['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useQuery<GetRenewalsQuery, TError, TData>({
    queryKey: ['getRenewals', variables],
    queryFn: fetcher<GetRenewalsQuery, GetRenewalsQueryVariables>(
      client,
      GetRenewalsDocument,
      variables,
      headers,
    ),
    ...options,
  });
};

useGetRenewalsQuery.document = GetRenewalsDocument;

useGetRenewalsQuery.getKey = (variables: GetRenewalsQueryVariables) => [
  'getRenewals',
  variables,
];

export const useInfiniteGetRenewalsQuery = <
  TData = InfiniteData<GetRenewalsQuery>,
  TError = unknown,
>(
  client: GraphQLClient,
  variables: GetRenewalsQueryVariables,
  options: Omit<
    UseInfiniteQueryOptions<GetRenewalsQuery, TError, TData>,
    'queryKey'
  > & {
    queryKey?: UseInfiniteQueryOptions<
      GetRenewalsQuery,
      TError,
      TData
    >['queryKey'];
  },
  headers?: RequestInit['headers'],
) => {
  return useInfiniteQuery<GetRenewalsQuery, TError, TData>(
    (() => {
      const { queryKey: optionsQueryKey, ...restOptions } = options;
      return {
        queryKey: optionsQueryKey ?? ['getRenewals.infinite', variables],
        queryFn: (metaData) =>
          fetcher<GetRenewalsQuery, GetRenewalsQueryVariables>(
            client,
            GetRenewalsDocument,
            { ...variables, ...(metaData.pageParam ?? {}) },
            headers,
          )(),
        ...restOptions,
      };
    })(),
  );
};

useInfiniteGetRenewalsQuery.getKey = (variables: GetRenewalsQueryVariables) => [
  'getRenewals.infinite',
  variables,
];

useGetRenewalsQuery.fetcher = (
  client: GraphQLClient,
  variables: GetRenewalsQueryVariables,
  headers?: RequestInit['headers'],
) =>
  fetcher<GetRenewalsQuery, GetRenewalsQueryVariables>(
    client,
    GetRenewalsDocument,
    variables,
    headers,
  );

useGetRenewalsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetRenewalsQueryVariables) =>
  (mutator: (cacheEntry: GetRenewalsQuery) => GetRenewalsQuery) => {
    const cacheKey = useGetRenewalsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<GetRenewalsQuery>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<GetRenewalsQuery>(cacheKey, mutator);
    }
    return { previousEntries };
  };
useInfiniteGetRenewalsQuery.mutateCacheEntry =
  (queryClient: QueryClient, variables?: GetRenewalsQueryVariables) =>
  (
    mutator: (
      cacheEntry: InfiniteData<GetRenewalsQuery>,
    ) => InfiniteData<GetRenewalsQuery>,
  ) => {
    const cacheKey = useInfiniteGetRenewalsQuery.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<InfiniteData<GetRenewalsQuery>>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<InfiniteData<GetRenewalsQuery>>(
        cacheKey,
        mutator,
      );
    }
    return { previousEntries };
  };
