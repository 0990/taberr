using System;
using System.IO;

namespace csharptest
{
    class Program
    {

        static void Main(string[] args)
        {
            using (var stream = new FileStream("../../../../Config.bin", FileMode.Open))
            {
                stream.Position = 0;

                var reader = new tabtoy.DataReader(stream);
                
                if ( !reader.ReadHeader( ) )
                {
                    Console.WriteLine("combine file crack!");
                    return;
                }

                var config = new gamedef.Config();
                config.Deserialize(reader);

                // ֱ��ͨ���±��ȡ�����
                var directFetch = config.Sample[2];

                // ��������ȡ
                var indexFetch = config.GetSampleByID(100);

                // ȡ�����ڵ�Ԫ��ʱ, ���ظ�����Ĭ��ֵ, �����
                var indexFetchByDefault = config.GetSampleByID(0, new gamedef.SampleDefine() );

                // �����־������Զ������
                config.TableLogger.AddTarget( new tabtoy.DebuggerTarget() );

                // ȡ��ʱ, ��Ĭ��ֵ��Ϊ��ʱ, �����־
                var nullFetchOutLog = config.GetSampleByID( 0 );

            }
            
            
            
        }

    }
}
