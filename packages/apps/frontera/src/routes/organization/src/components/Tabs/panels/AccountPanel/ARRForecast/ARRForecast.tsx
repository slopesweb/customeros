import { useParams } from 'react-router-dom';

import { useIsRestoring } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { useStore } from '@shared/hooks/useStore';
import { HelpCircle } from '@ui/media/icons/HelpCircle';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { Card, CardContent } from '@ui/presentation/Card/Card';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog/InfoDialog';
import {
  Contract,
  RenewalSummary,
  OpportunityRenewalLikelihood,
} from '@graphql/types';
import { getRenewalLikelihoodColor } from '@organization/components/Tabs/panels/AccountPanel/utils';
import { useARRInfoModalContext } from '@organization/components/Tabs/panels/AccountPanel/context/AccountModalsContext';

interface ARRForecastProps {
  name: string;
  currency?: string | null;
  isInitialLoading?: boolean;
  contracts?: Contract[] | null;
  renewalSunnary?: RenewalSummary | null;
}

export const ARRForecast = ({
  isInitialLoading,
  renewalSunnary,
  name,
  currency = 'USD',
}: ARRForecastProps) => {
  const isRestoring = useIsRestoring();
  const store = useStore();
  const id = useParams()?.id as string;

  const organization = store.organizations.value.get(id);

  const { modal } = useARRInfoModalContext();
  const formattedMaxAmount = formatCurrency(
    renewalSunnary?.maxArrForecast ?? 0,
    2,
    currency || 'USD',
  );
  const formattedAmount = formatCurrency(
    renewalSunnary?.arrForecast ?? 0,
    2,
    currency || 'USD',
  );

  const hasForecastChanged = formattedMaxAmount !== formattedAmount;

  return (
    <>
      <Card className='p-4 w-full bg-transparent cursor-default group border-0'>
        <CardContent className='p-0 flex items-center '>
          <FeaturedIcon
            size='md'
            colorScheme={getRenewalLikelihoodColor(
              renewalSunnary?.renewalLikelihood,
            )}
            className={
              renewalSunnary?.renewalLikelihood ===
              OpportunityRenewalLikelihood.LowRenewal
                ? 'text-orangeDark-800'
                : undefined
            }
          >
            <CurrencyDollar />
          </FeaturedIcon>
          <div className='flex ml-5 w-full items-center gap-4 justify-between'>
            <div className='flex flex-col'>
              <div className='flex items-center'>
                <h2 className='whitespace-nowrap font-semibold text-gray-700 mr-2'>
                  ARR forecast
                </h2>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='Help'
                  icon={<HelpCircle className='text-gray-400' />}
                  className='group-hover:opacity-100 opacity-0 transition-opacity duration-200 ease-linear'
                  onClick={(e) => {
                    e.stopPropagation();
                    modal.onOpen();
                  }}
                />
              </div>
            </div>

            <div className='flex flex-col'>
              <h2
                className={cn(
                  organization?.isLoading &&
                    (!isInitialLoading || !isRestoring
                      ? 'text-gray-400'
                      : 'text-gray-700'),
                  'text-2xl font-semibold transition-opacity duration-250 ease-in',
                )}
              >
                {formattedAmount}
              </h2>
              {hasForecastChanged && !organization?.isLoading && (
                <p className='text-sm  text-right line-through'>
                  {formattedMaxAmount}
                </p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      <InfoDialog
        isOpen={modal.open}
        label='ARR forecast'
        onClose={modal.onClose}
        onConfirm={modal.onClose}
        confirmButtonLabel='Got it'
      >
        <p className='text-sm mb-4 text-gray-700'>
          Annual Recurring Revenue (ARR) is the total amount of money you can
          expect to receive from
          <span className='font-medium mx-1'>{name ? name : `Unnamed`}</span>
          for the next 12 months.
        </p>
        <p className='text-sm font-normal text-gray-700'>
          It includes all renewals but excludes one-time and per use services.
          Renewals are discounted based on the renewal likelihood
        </p>
      </InfoDialog>
    </>
  );
};
