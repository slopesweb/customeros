import React, { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { cn } from '@ui/utils/cn.ts';
import { InputProps } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { SelectOption } from '@shared/types/SelectOptions.ts';
import { Divider } from '@ui/presentation/Divider/Divider.tsx';
import { validateEmail } from '@shared/util/emailValidation.ts';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';
import { EmailSelect } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/components/EmailsInputGroup/EmailSelect.tsx';

interface EmailsInputGroupProps extends InputProps {
  contractId: string;
}

const EmailList = ({ emailList }: { emailList: string[] }) => {
  return (
    <p className='text-gray-500 whitespace-nowrap overflow-ellipsis overflow-hidden h-8 mt-1 border-b-1 border-transparent flex'>
      {[...emailList].map((email, i) => {
        const validationMessage = validateEmail(email);

        return (
          <React.Fragment key={email}>
            <Tooltip label={validationMessage || ''}>
              <span
                className={cn('mr-1 text-base my-1 block', {
                  'text-warning-700': validateEmail(email),
                })}
              >
                {email}
                {i < emailList.length - 1 && ','}
              </span>
            </Tooltip>
          </React.Fragment>
        );
      })}
    </p>
  );
};

export const EmailsInputGroup = observer(
  ({ contractId }: EmailsInputGroupProps) => {
    const store = useStore();

    const contractStore = store.contracts.value.get(
      contractId,
    ) as ContractStore;
    const billingDetails = contractStore?.tempValue?.billingDetails;

    const [showCC, setShowCC] = useState(false);
    const [showBCC, setShowBCC] = useState(false);
    const [isFocused, setIsFocused] = useState(false);
    const [focusedItemIndex, setFocusedItemIndex] = useState<false | number>(
      false,
    );
    const ref = React.useRef(null);

    useOutsideClick({
      ref: ref,
      handler: () => {
        setIsFocused(false);
        setFocusedItemIndex(false);
        setShowCC(false);
        setShowBCC(false);
      },
    });

    const handleFocus = (index: number) => {
      setIsFocused(true);
      setFocusedItemIndex(index);
    };

    useEffect(() => {
      if (showCC && !isFocused) {
        handleFocus(1);
      }
    }, [showCC]);

    useEffect(() => {
      if (showBCC && !isFocused) {
        handleFocus(2);
      }
    }, [showBCC]);

    const handleUpdateBillingEmailData = (
      key: string,
      value: SelectOption<string> | SelectOption<string>[],
    ) => {
      const val = Array.isArray(value)
        ? value.map((v) => v.value)
        : value?.value;

      contractStore?.updateTemp((contract) => ({
        ...contract,
        billingDetails: {
          ...contract.billingDetails,
          [key]: val,
        },
      }));
    };

    return (
      <div ref={ref}>
        <div className='flex relative items-center h-8  mt-2'>
          <p className='text-sm text-gray-500 after:border-t-2 w-fit whitespace-nowrap mr-2'>
            Send invoice
          </p>
          <Divider />

          <div className='flex'>
            {!showCC && !billingDetails?.billingEmailCC?.length && (
              <Button
                size='sm'
                variant='ghost'
                color='gray.400'
                className='text-sm px-1 mx-1'
                onClick={() => {
                  setShowCC(true);
                  setFocusedItemIndex(1);
                }}
              >
                CC
              </Button>
            )}

            {!showBCC && !billingDetails?.billingEmailBCC?.length && (
              <Button
                size='sm'
                variant='ghost'
                color='gray.400'
                className='text-sm px-1 '
                onClick={() => {
                  setShowBCC(true);
                  setFocusedItemIndex(2);
                }}
              >
                BCC
              </Button>
            )}
          </div>
        </div>
        <EmailSelect
          entryType='To'
          isMulti={false}
          placeholder='To email address'
          autofocus={focusedItemIndex === 0}
          value={
            billingDetails?.billingEmail ? [billingDetails.billingEmail] : []
          }
          onChange={(value: SelectOption<string>[]) =>
            handleUpdateBillingEmailData('billingEmail', value)
          }
        />

        <div className='flex-col flex-1 w-full gap-4'>
          {isFocused && (
            <>
              {(showCC || !!billingDetails?.billingEmailCC?.length) && (
                <EmailSelect
                  isMulti
                  entryType='CC'
                  placeholder='CC email addresses'
                  autofocus={focusedItemIndex === 1}
                  value={billingDetails?.billingEmailCC ?? []}
                  onChange={(value: SelectOption<string>[]) =>
                    handleUpdateBillingEmailData('billingEmailCC', value)
                  }
                />
              )}
              {(showBCC || !!billingDetails?.billingEmailBCC?.length) && (
                <EmailSelect
                  isMulti
                  entryType='BCC'
                  placeholder='BCC email addresses'
                  autofocus={focusedItemIndex === 2}
                  value={billingDetails?.billingEmailBCC ?? []}
                  onChange={(value: SelectOption<string>[]) =>
                    handleUpdateBillingEmailData('billingEmailBCC', value)
                  }
                />
              )}
            </>
          )}
        </div>

        {!isFocused && (
          <div className='flex-col flex-1 gap-4'>
            {!!billingDetails?.billingEmailCC?.length && (
              <div
                role='button'
                onClick={() => handleFocus(1)}
                aria-label='Click to input participant data'
                className={cn('overflow-hidden', {
                  'flex-1': !billingDetails?.billingEmailBCC?.length,
                })}
              >
                <span className='text-sm font-semibold text-gray-700 mr-1'>
                  CC
                </span>
                <EmailList emailList={billingDetails?.billingEmailCC ?? []} />
              </div>
            )}
            {!!billingDetails?.billingEmailBCC?.length && (
              <div
                role='button'
                onClick={() => handleFocus(2)}
                aria-label='Click to input participant data'
                className={cn('overflow-hidden', {
                  'flex-1': !billingDetails?.billingEmailBCC?.length,
                })}
              >
                <span className='text-sm font-semibold text-gray-700 mr-1'>
                  BCC
                </span>
                <EmailList emailList={billingDetails?.billingEmailBCC ?? []} />
              </div>
            )}
          </div>
        )}
      </div>
    );
  },
);
