using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading.Tasks;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using TwoMQTT.Core;
using TwoMQTT.Core.Extensions;
using TwoMQTT.Core.Interfaces;
using TwoMQTT.Core.Utils;
using WSDOT.DataAccess;
using WSDOT.Liasons;
using WSDOT.Models.Shared;

namespace WSDOT
{
    class Program : ConsoleProgram<Resource, Command, SourceLiason, MQTTLiason>
    {
        static async Task Main(string[] args)
        {
            var p = new Program();
            await p.ExecuteAsync(args);
        }

        /// <inheritdoc />
        protected override IDictionary<string, string> EnvironmentDefaults()
        {
            var sep = "__";
            var section = Models.Options.MQTTOpts.Section.Replace(":", sep);
            var sectsep = $"{section}{sep}";

            return new Dictionary<string, string>
            {
                { $"{sectsep}{nameof(Models.Options.MQTTOpts.TopicPrefix)}", Models.Options.MQTTOpts.TopicPrefixDefault },
                { $"{sectsep}{nameof(Models.Options.MQTTOpts.DiscoveryName)}", Models.Options.MQTTOpts.DiscoveryNameDefault },
            };
        }

        /// <inheritdoc />
        protected override IServiceCollection ConfigureServices(HostBuilderContext hostContext, IServiceCollection services)
        {
            services.AddHttpClient<ISourceDAO>();

            return services
                .ConfigureOpts<Models.Options.SharedOpts>(hostContext, Models.Options.SharedOpts.Section)
                .ConfigureOpts<Models.Options.SourceOpts>(hostContext, Models.Options.SourceOpts.Section)
                .ConfigureOpts<TwoMQTT.Core.Models.MQTTManagerOptions>(hostContext, Models.Options.MQTTOpts.Section)
                .AddSingleton<IThrottleManager, ThrottleManager>(x =>
                {
                    var opts = x.GetService<IOptions<Models.Options.SourceOpts>>();
                    return new ThrottleManager(opts.Value.PollingInterval);
                })
                .AddSingleton<ISourceDAO, SourceDAO>(x =>
                {
                    var opts = x.GetService<IOptions<Models.Options.SourceOpts>>();
                    return new SourceDAO(x.GetService<ILogger<SourceDAO>>(), x.GetService<IHttpClientFactory>(), opts.Value.ApiKey);
                });
        }
    }
}