import { useSearchParams } from 'react-router-dom';
import { useRef, useEffect, startTransition } from 'react';

import { match } from 'ts-pattern';
import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Input } from '@ui/form/Input/Input';
import { Star06 } from '@ui/media/icons/Star06';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { TableIdType, TableViewType } from '@graphql/types';
import { ViewSettings } from '@shared/components/ViewSettings';
import { UserPresence } from '@shared/components/UserPresence';
import { TableViewMenu } from '@organizations/components/TableViewMenu/TableViewMenu';
import {
  InputGroup,
  LeftElement,
  RightElement,
} from '@ui/form/InputGroup/InputGroup';
import { TableViewsToggleNavigation } from '@organizations/components/TableViewsToggleNavigation';
import { SearchBarFilterData } from '@organizations/components/SearchBarFilterData/SearchBarFilterData';

interface SearchProps {
  open: boolean;
  onOpen: () => void;
  onClose: () => void;
}

export const Search = observer(({ onClose, onOpen, open }: SearchProps) => {
  const store = useStore();
  const wrapperRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const measureRef = useRef<HTMLDivElement>(null);
  const floatingActionPropmterRef = useRef<HTMLDivElement>(null);
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const [lastSearchForPreset] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-last-search-for-preset`, { root: 'root' });

  const displayIcp = useFeatureIsOn('icp');

  useEffect(() => {
    onClose();

    setSearchParams(
      (prev) => {
        if (preset && lastSearchForPreset?.[preset]) {
          prev.set('search', lastSearchForPreset[preset]);
        } else {
          prev.delete('search');
        }

        return prev;
      },
      { replace: true },
    );

    if (preset && inputRef?.current) {
      inputRef.current.value = lastSearchForPreset[preset] ?? '';
    }
  }, [preset]);

  const tableViewType = store.tableViewDefs.getById(preset || '')?.value
    .tableType;
  const tableId = store.tableViewDefs.getById(preset || '')?.value.tableId;

  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableType = tableViewDef?.value?.tableType;
  const totalResults = store.ui.searchCount;

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    startTransition(() => {
      const value = event.target.value;

      setSearchParams(
        (prev) => {
          if (!value) {
            prev.delete('search');
          } else {
            prev.set('search', value);
          }

          return prev;
        },
        { replace: true },
      );
    });
  };

  useKeyBindings(
    {
      '/': () => {
        setTimeout(() => {
          inputRef.current?.focus();
        }, 0);
      },
    },
    {
      when:
        !store.ui.isEditingTableCell &&
        !store.ui.isFilteringTable &&
        !store.ui.commandMenu.isOpen,
    },
  );

  const placeholder = match(tableType)
    .with(TableViewType.Contacts, () => 'e.g. Isabella Evans')
    .with(TableViewType.Contracts, () => 'e.g. CustomerOS contract...')
    .with(TableViewType.Organizations, () => 'e.g. CustomerOS...')
    .with(TableViewType.Invoices, () => 'e.g. My contract')
    .otherwise(() => 'e.g. Organization name...');

  const handleToogleFlow = () => {
    if (open) {
      onClose();
    } else {
      onOpen();
    }
  };

  const allowCreation =
    tableViewDef?.value.tableId !== TableIdType.Contracts &&
    totalResults === 0 &&
    !!searchParams.get('search');

  useKeyBindings(
    {
      Enter: () => {
        store.ui.commandMenu.setType('AddNewOrganization');
        store.ui.commandMenu.setOpen(true);
      },
    },
    { when: allowCreation },
  );

  return (
    <div
      ref={wrapperRef}
      className='flex items-center justify-between pr-1 w-full data-[focused]:animate-focus gap-2'
    >
      <InputGroup className='relative w-full bg-transparent hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent gap-1'>
        <LeftElement className='ml-2'>
          <SearchBarFilterData />
        </LeftElement>
        <Input
          size='md'
          ref={inputRef}
          autoCorrect='off'
          spellCheck={false}
          variant='unstyled'
          onChange={handleChange}
          defaultValue={searchParams.get('search') ?? ''}
          placeholder={
            store.ui.isSearching !== 'organizations'
              ? `/ to search`
              : placeholder
          }
          onBlur={() => {
            store.ui.setIsSearching(null);
            wrapperRef.current?.removeAttribute('data-focused');
          }}
          onFocus={() => {
            store.ui.setIsSearching('organizations');
            wrapperRef.current?.setAttribute('data-focused', '');
          }}
          onKeyDown={(e) => {
            if (e.key === 'a' && (e.metaKey || e.ctrlKey)) {
              e.preventDefault();
              e.currentTarget.select();
            }
          }}
          onKeyUp={(e) => {
            if (
              e.code === 'Escape' ||
              e.code === 'ArrowUp' ||
              e.code === 'ArrowDown'
            ) {
              inputRef.current?.blur();
              store.ui.setIsSearching(null);
            }
          }}
        />
        <RightElement>
          {allowCreation && (
            <div
              role='button'
              ref={floatingActionPropmterRef}
              onClick={() => inputRef.current?.focus()}
              className='flex flex-row items-center gap-1 absolute top-[11px] cursor-text'
              style={{
                left: `calc(${measureRef?.current?.offsetWidth ?? 0}px + 24px)`,
              }}
            >
              <Tag variant='subtle' className='mb-[2px]' colorScheme='grayBlue'>
                <TagLabel className='capitalize'>Enter</TagLabel>
              </Tag>
              <span className='font-normal text-gray-400 break-keep w-max text-sm'>
                to create
              </span>
            </div>
          )}
        </RightElement>
      </InputGroup>
      <UserPresence channelName={`finder:${store.session.value.tenant}`} />

      <TableViewsToggleNavigation />

      {tableViewType && <ViewSettings type={tableViewType} />}

      {TableIdType.Nurture === tableId && displayIcp && (
        <IconButton
          size='xs'
          icon={<Star06 />}
          aria-label='toogle-flow'
          onClick={handleToogleFlow}
        />
      )}
      <TableViewMenu />
      <span ref={measureRef} className={`z-[-1] absolute h-0 invisible flex`}>
        <div className='ml-2'>
          <SearchBarFilterData />
        </div>
        {inputRef?.current?.value ?? ''}
      </span>
    </div>
  );
});
