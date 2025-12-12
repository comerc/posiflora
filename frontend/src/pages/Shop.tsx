import {
  Button,
  Card,
  Form,
  Input,
  Switch,
  Checkbox,
  message,
  Statistic,
  Divider,
  Alert,
} from 'antd'
import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { format } from 'date-fns'
import { api, type ConnectTelegramRequest, type TelegramStatus } from '@/services/api'

function ShopPage() {
  const { shopId } = useParams<{ shopId: string }>()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [statusLoading, setStatusLoading] = useState(false)
  const [status, setStatus] = useState<TelegramStatus | null>(null)
  const [error, setError] = useState<string | null>(null)
  // Начальное состояние: галочка не нажата (false), нужно маскировать
  // В бэкенде maskDisabled=true маскирует, поэтому передаем !showFullChatId
  const [showFullChatId, setShowFullChatId] = useState(false)

  // Загрузка статуса
  const loadStatus = async (showError = false, silent = false, showFullChatId = false) => {
    if (!shopId) return

    if (!silent) {
      setStatusLoading(true)
    }
    setError(null)
    try {
      // В бэкенде логика инвертирована: maskDisabled=true маскирует, maskDisabled=false показывает полный
      // showFullChatId=true означает "показать полный Chat ID" → передаем maskDisabled=false
      // showFullChatId=false означает "маскировать Chat ID" → передаем maskDisabled=true
      // Инвертируем, чтобы соответствовать инвертированной логике бэкенда

      const data = await api.getTelegramStatus(Number(shopId), showFullChatId)
      setStatus(data)
      // Обновляем только enabled в форме (токены не показываем из соображений безопасности)
      form.setFieldsValue({
        enabled: data.enabled,
      })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Не удалось загрузить статус'
      setError(errorMessage)
      // Показываем ошибку только при первой загрузке или если явно запрошено
      if (showError || status === null) {
        message.error(errorMessage)
      }
    } finally {
      if (!silent) {
        setStatusLoading(false)
      }
    }
  }

  useEffect(() => {
    if (!shopId) return
    loadStatus()
  }, [shopId])

  // Сохранение интеграции
  const handleSubmit = async (values: ConnectTelegramRequest) => {
    if (!shopId) return

    setLoading(true)
    setError(null)
    try {
      await api.connectTelegram(Number(shopId), values)
      message.success('Интеграция успешно сохранена')
      // Обновляем статус после успешного сохранения (показываем ошибку если не удалось)
      await loadStatus(true)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Не удалось сохранить интеграцию'
      message.error(errorMessage)
      setError(errorMessage)
      console.error('Ошибка при сохранении интеграции:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container mx-auto max-w-6xl p-6">
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Левая колонка - форма и статус */}
        <div className="space-y-6">
          <Card title="Telegram интеграция">
            <Form
              autoComplete="off"
              form={form}
              layout="vertical"
              onFinish={handleSubmit}
              initialValues={{
                enabled: false,
              }}
            >
              <Form.Item
                label="Bot Token"
                name="botToken"
                rules={[{ required: true, message: 'Введите Bot Token' }]}
              >
                <Input.Password
                  placeholder="123456:ABC-DEF..."
                  size="large"
                  autoComplete="new-password"
                  name="bot-token-hidden"
                  id="bot-token-hidden"
                />
              </Form.Item>

              <Form.Item
                label="Chat ID"
                name="chatId"
                rules={[{ required: true, message: 'Введите Chat ID' }]}
              >
                <Input
                  placeholder="987654321"
                  size="large"
                  autoComplete="off"
                  name="chat-id-hidden"
                  id="chat-id-hidden"
                />
              </Form.Item>

              <Form.Item label="Включено" name="enabled" valuePropName="checked">
                <Switch />
              </Form.Item>

              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading} size="large" block>
                  Сохранить
                </Button>
              </Form.Item>
            </Form>
          </Card>

          <Card
            title="Статус интеграции"
            loading={statusLoading}
            extra={
              <Checkbox
                checked={showFullChatId}
                onChange={(e) => {
                  const newValue = e.target.checked
                  setShowFullChatId(newValue)
                  // Перезагружаем статус без показа loading, чтобы не было мигания
                  loadStatus(false, true, newValue)
                }}
              >
                Показать Chat ID
              </Checkbox>
            }
          >
            {error && !status ? (
              <div className="py-4 text-center text-red-500">{error}</div>
            ) : status ? (
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <Statistic
                    title="Статус"
                    value={status.enabled ? 'Включена' : 'Выключена'}
                    valueStyle={{ color: status.enabled ? '#3f8600' : '#cf1322' }}
                  />
                  <Statistic
                    groupSeparator=""
                    title="Chat ID"
                    value={status.chat_id || 'Не настроен'}
                  />
                </div>

                {status.last_sent_at && (
                  <div>
                    <Statistic
                      title="Последняя отправка"
                      value={format(new Date(status.last_sent_at), 'dd.MM.yyyy HH:mm:ss')}
                    />
                  </div>
                )}

                <Divider />

                <div className="grid grid-cols-2 gap-4">
                  <Statistic
                    title="Отправлено за 7 дней"
                    value={status.sent_count}
                    valueStyle={{ color: '#3f8600' }}
                  />
                  <Statistic
                    title="Ошибок за 7 дней"
                    value={status.failed_count}
                    valueStyle={{ color: '#cf1322' }}
                  />
                </div>
              </div>
            ) : (
              <div className="py-4 text-center text-gray-500">Интеграция не настроена</div>
            )}
          </Card>
        </div>

        {/* Правая колонка - инструкция */}
        <div>
          <Card title="Как узнать Chat ID?">
            <Alert
              message="Инструкция"
              description={
                <div className="mt-2 space-y-2">
                  <p>
                    1. Создайте бота через{' '}
                    <a
                      href="https://t.me/BotFather"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-blue-500"
                    >
                      @BotFather
                    </a>{' '}
                    в Telegram
                  </p>
                  <p>2. Получите Bot Token от @BotFather</p>
                  <p>3. Для получения Chat ID:</p>
                  <ul className="ml-4 list-inside list-disc space-y-1">
                    <li>
                      Напишите боту{' '}
                      <a
                        href="https://t.me/userinfobot"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-blue-500"
                      >
                        @userinfobot
                      </a>{' '}
                      - он покажет ваш Chat ID
                    </li>
                    <li>Или добавьте бота в группу и используйте ID группы</li>
                    <li>
                      Или используйте{' '}
                      <a
                        href="https://t.me/getidsbot"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-blue-500"
                      >
                        @getidsbot
                      </a>
                    </li>
                  </ul>
                </div>
              }
              type="info"
              showIcon
            />
          </Card>
        </div>
      </div>
    </div>
  )
}

export default ShopPage
