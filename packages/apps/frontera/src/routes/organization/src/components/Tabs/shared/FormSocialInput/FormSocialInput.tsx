import { useField } from 'react-inverted-form';
import { useRef, useMemo, useState, useCallback } from 'react';

import { Input } from '@ui/form/Input';
import { Social } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useUpdateSocialMutation } from '@organization/graphql/updateSocial.generated';
import { useRemoveSocialMutation } from '@organization/graphql/removeSocial.generated';
import {
  InputGroup,
  LeftElement,
  InputGroupProps,
} from '@ui/form/InputGroup/InputGroup';

import { SocialIcon } from './SocialIcons';
import { SocialInput } from './SocialInput';

type SocialInputValue = Pick<Social, 'id' | 'url'>;

interface FormSocialInputProps extends InputGroupProps {
  name: string;
  formId: string;
  isReadOnly?: boolean;
  placeholder?: string;
  organizationId: string;
  invalidateQuery: () => void;
  leftElement?: React.ReactNode;
  defaultValues: Array<SocialInputValue>;
  addSocial: (props: {
    newValue: string;
    onSuccess: ({ id, url }: { id: string; url: string }) => void;
  }) => void;
}

export const FormSocialInput = ({
  name,
  formId,
  leftElement,
  isReadOnly,
  organizationId,
  defaultValues,
  addSocial,
  invalidateQuery,
  ...rest
}: FormSocialInputProps) => {
  const { getInputProps } = useField(name, formId);
  const { value, onChange, onBlur } = getInputProps();
  const values = useMemo(
    () => (Array.isArray(value) ? ([...value] as SocialInputValue[]) : value),
    [value],
  );
  const _leftElement = useMemo(() => leftElement, []);

  const client = getGraphQLClient();

  const updateSocial = useUpdateSocialMutation(client, {
    onSuccess: invalidateQuery,
  });
  const removeSocial = useRemoveSocialMutation(client, {
    onSuccess: invalidateQuery,
  });

  const newInputRef = useRef<HTMLInputElement>(null);
  const [newValue, setNewValue] = useState('');

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const id = e?.target?.id;
      const next = [...values];
      const index = next.findIndex((item) => item.id === id);

      next[index].url = e.target.value?.trim();
      onChange(next);
    },
    [values],
  );

  const handleBlur = useCallback(
    (e: React.FocusEvent<HTMLInputElement>) => {
      const next = [...values];
      const index = next.findIndex((item) => item.id === e.target.id);

      if (!e.target.value) {
        removeSocial.mutate(
          { socialId: values[index].id },
          {
            onSuccess: () => {
              next.splice(index, 1);
              onBlur?.(next);
            },
          },
        );
      } else {
        const { id, url } = values[index];

        updateSocial.mutate(
          { input: { id, url } },
          {
            onSuccess: () => {
              onBlur?.(values);
            },
          },
        );
      }
    },
    [values],
  );

  const handleRemoveKeyDown = useCallback(
    (e: React.KeyboardEvent<HTMLInputElement>) => {
      const next = [...values];
      const index = next.findIndex((item) => item.id === e.currentTarget.id);

      if (e.key === 'Backspace' && !values[index].url) {
        removeSocial.mutate(
          { socialId: values[index].id },
          {
            onSuccess: () => {
              next.splice(index, 1);
              onBlur?.(next);
              newInputRef.current?.focus();
            },
          },
        );
      }
    },
    [values],
  );

  const handleAddKeyDown = useCallback(
    (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Enter') {
        if (newValue) {
          addSocial({
            newValue,
            onSuccess: ({ id, url }: { id: string; url: string }) => {
              onBlur?.([...values, { id, url }]);
              setNewValue('');
            },
          });
        }
      }
    },
    [newValue, organizationId, values],
  );

  const handleAddChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setNewValue(e.target.value);
    },
    [],
  );

  const handleAddBlur = useCallback(() => {
    if (newValue) {
      addSocial({
        newValue,
        onSuccess: ({ id, url }: { id: string; url: string }) => {
          onBlur?.([...values, { id, url }]);
          setNewValue('');
        },
      });
    }
  }, [newValue, organizationId, values]);

  return (
    <>
      {((values as SocialInputValue[]) || [])?.map(({ id, url }, index) => (
        <SocialInput
          id={id}
          key={index}
          value={url}
          index={index}
          onBlur={handleBlur}
          isReadOnly={isReadOnly}
          onChange={handleChange}
          leftElement={_leftElement}
          onKeyDown={handleRemoveKeyDown}
        />
      ))}

      {!isReadOnly && (
        <InputGroup {...rest}>
          {leftElement && (
            <LeftElement>
              <SocialIcon url={newValue}>{leftElement}</SocialIcon>
            </LeftElement>
          )}
          <Input
            value={newValue}
            ref={newInputRef}
            onBlur={handleAddBlur}
            onChange={handleAddChange}
            onKeyDown={handleAddKeyDown}
            className={
              'border-b border-transparent hover:border-transparent hover:border-b-none text-md focus:hover:border-b focus:hover:border-transparent focus:border-b focus:border-transparent'
            }
            {...rest}
          />
        </InputGroup>
      )}
    </>
  );
};
